package scheduler

import (
	"encoding/json"
	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/models"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

var (
	c        *cron.Cron
	handlers map[string]func()
	mu       sync.Mutex
	App      *fiber.App
)

// Init initializes the scheduler
func Init(app *fiber.App) {
	App = app
	// Initialize standard cron parser
	c = cron.New()
	handlers = make(map[string]func())
	c.Start()
	log.Println("Scheduler started")
}

// RegisterHandler registers a function to be called by a job name
func RegisterHandler(name string, handler func()) {
	mu.Lock()
	defer mu.Unlock()
	handlers[name] = handler
	log.Printf("Registered job handler: %s", name)
}

// LoadJobs loads jobs from database and schedules them
func LoadJobs(db *gorm.DB) {
	// Sync jobs from workflows first
	SyncJobsFromWorkflows(db)

	// Stop existing cron scheduler and create a new one to avoid duplicates
	mu.Lock()
	if c != nil {
		c.Stop()
	}
	c = cron.New()
	c.Start()
	mu.Unlock()
	log.Println("Cron scheduler reset and restarted")

	var jobs []models.Jobs

	// recordstatus = 1 means the job is active/enabled in the schedule
	if err := db.Where("recordstatus = ?", 1).Find(&jobs).Error; err != nil {
		log.Println("Error loading jobs:", err)
		return
	}

	log.Printf("Found %d active jobs in database", len(jobs))

	for _, job := range jobs {
		if job.Schedule == "" {
			log.Printf("Job %s has empty schedule, skipping", job.JobName)
			continue
		}

		// Capture job for closure
		j := job

		_, err := c.AddFunc(j.Schedule, func() {
			runJob(db, j)
		})

		if err != nil {
			log.Printf("Error scheduling job %s with schedule '%s': %v\n", j.JobName, j.Schedule, err)
		} else {
			log.Printf("Scheduled job: %s [%s]\n", j.JobName, j.Schedule)
		}
	}
}

func runJob(gormDB *gorm.DB, job models.Jobs) {
	mu.Lock()
	handler, exists := handlers[job.JobName]
	mu.Unlock()

	// If Flow is set, we don't strictly *need* a handler, but we can have both or either.
	// Logic: If Flow is present, execute flow. If handler is present, execute handler.
	// If neither, warn.

	if !exists && job.Flow == "" {
		log.Printf("No handler registered and no Flow defined for job: %s\n", job.JobName)
		return
	}

	log.Printf("Starting job: %s\n", job.JobName)

	// Set isrunning = 1 (Running)
	if err := gormDB.Model(&models.Jobs{}).Where("jobsid = ?", job.JobsID).Update("isrunning", 1).Error; err != nil {
		log.Printf("Failed to set isrunning=1 for job %s: %v", job.JobName, err)
	}

	// Execute handler
	// Recover from panics in handlers to prevent crashing the scheduler/app
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Job %s panicked: %v", job.JobName, r)
		}

		// Update LastRunning and Set isrunning = 0 (Finished)
		updates := map[string]interface{}{
			"isrunning":   0,
			"lastrunning": time.Now(),
		}
		if err := gormDB.Model(&models.Jobs{}).Where("jobsid = ?", job.JobsID).Updates(updates).Error; err != nil {
			log.Printf("Failed to update status for job %s: %v", job.JobName, err)
		}
		log.Printf("Finished job: %s\n", job.JobName)
	}()

	if exists {
		handler()
	}

	// Execute external command if specified
	var commandOutput string
	if job.Executable != "" {
		log.Printf("Executing command: %s", job.Executable)
		output, err := executeCommand(job.Executable)
		commandOutput = output
		if err != nil {
			log.Printf("Command execution failed: %v", err)
			commandOutput = fmt.Sprintf("Error: %v\nOutput: %s", err, output)
		} else {
			log.Printf("Command executed successfully")
		}
	}

	if job.Flow != "" {
		if App == nil {
			log.Println("Fiber App not initialized in scheduler, cannot execute flow")
			return
		}
		log.Printf("Executing Flow: %s", job.Flow)

		// Create a temporary context for the flow execution
		ctx := App.AcquireCtx(&fasthttp.RequestCtx{})
		defer App.ReleaseCtx(ctx)

		// We need to initialize locals that ExecuteFlow expects
		// ExecuteFlow uses "wfEngine", "flowTerminated"
		var wfEngine []generator.WorkflowEngine

		// Store command output in context so components can access it
		if commandOutput != "" {
			ctx.Locals("commandOutput", commandOutput)
			log.Printf("Stored command output in context: %s", commandOutput)
		}

		ctx.Locals("wfEngine", wfEngine)
		ctx.Locals("flowTerminated", false)

		if err := generator.ExecuteFlow(ctx, gormDB, job.Flow, false); err != nil {
			log.Printf("Error executing flow %s: %v", job.Flow, err)
		} else {
			log.Printf("Flow %s executed successfully", job.Flow)
		}
	}
}

// SyncJobsFromWorkflows scans active workflows for the 'schedule' component and updates the jobs table.
func SyncJobsFromWorkflows(db *gorm.DB) {
	log.Println("SyncJobsFromWorkflows: Starting workflow sync...")

	var workflows []models.Workflow
	if err := db.Where("recordstatus = 1").Find(&workflows).Error; err != nil {
		log.Println("Error loading workflows for sync:", err)
		return
	}

	log.Printf("SyncJobsFromWorkflows: Found %d active workflows", len(workflows))

	for _, wf := range workflows {
		var flowData generator.FlowData
		if err := json.Unmarshal([]byte(wf.Flow), &flowData); err != nil {
			log.Printf("SyncJobsFromWorkflows: Failed to parse workflow %s: %v", wf.Wfname, err)
			continue
		}

		foundSchedule := false
		var cronExp string
		var enabled string
		var executable string
		var nodeID int

		// Iterate through components to find 'schedule' or 'scheduler'
		for _, comp := range flowData.Drawflow.Home.Data {
			compNameLower := strings.ToLower(comp.Name)
			if compNameLower == "schedule" || compNameLower == "scheduler" {
				foundSchedule = true
				nodeID = comp.ID
				log.Printf("SyncJobsFromWorkflows: Found scheduler component in workflow %s (nodeID: %d)", wf.Wfname, nodeID)
				break
			}
		}

		// If scheduler component found, get its properties from workflowdetail table
		if foundSchedule {
			var details []models.Workflowdetail
			if err := db.Preload("Componentdetail").Where("workflowid = ? AND nodeid = ?", wf.Workflowid, nodeID).Find(&details).Error; err != nil {
				log.Printf("SyncJobsFromWorkflows: Failed to load details for workflow %s node %d: %v", wf.Wfname, nodeID, err)
				continue
			}

			// Extract cron, enabled, and executable from workflowdetail
			for _, detail := range details {
				if strings.ToLower(detail.Componentdetail.Inputname) == "cron" {
					cronExp = detail.Componentvalue
				}
				if strings.ToLower(detail.Componentdetail.Inputname) == "enabled" {
					enabled = detail.Componentvalue
				}
				if strings.ToLower(detail.Componentdetail.Inputname) == "executable" {
					executable = detail.Componentvalue
				}
			}
			log.Printf("SyncJobsFromWorkflows: Workflow %s scheduler config - cron: %s, enabled: %s, executable: %s", wf.Wfname, cronExp, enabled, executable)
		}

		// Check if we need to foster a job for this workflow
		if foundSchedule && (enabled == "true" || enabled == "1" || strings.ToLower(enabled) == "true") && cronExp != "" {
			var job models.Jobs
			err := db.Where("jobname = ?", wf.Wfname).First(&job).Error

			if err == gorm.ErrRecordNotFound {
				// Create new job
				var maxID int
				db.Model(&models.Jobs{}).Select("IFNULL(MAX(jobsid), 0)").Scan(&maxID)

				newJob := models.Jobs{
					JobsID:       maxID + 1,
					JobName:      wf.Wfname,
					Description:  "Auto-generated from workflow " + wf.Wfname,
					Version:      "1.0",
					CreatedBy:    "system",
					Flow:         wf.Wfname,
					Executable:   executable,
					Schedule:     cronExp,
					RecordStatus: 1,
					LastRunning:  time.Now(),
				}

				if err := db.Create(&newJob).Error; err != nil {
					log.Printf("Failed to create auto-job for workflow %s: %v", wf.Wfname, err)
				} else {
					log.Printf("Created auto-job for workflow %s with schedule %s", wf.Wfname, cronExp)
				}
			} else if err == nil {
				// Update existing job if needed
				if job.Schedule != cronExp || job.RecordStatus != 1 || job.Flow != wf.Wfname {
					if err := db.Model(&job).Updates(map[string]interface{}{
						"schedule":     cronExp,
						"recordstatus": 1,
						"flow":         wf.Wfname,
						"executable":   executable,
					}).Error; err != nil {
						log.Printf("Failed to update auto-job for workflow %s: %v", wf.Wfname, err)
					} else {
						log.Printf("Updated auto-job for workflow %s with new schedule %s", wf.Wfname, cronExp)
					}
				}
			}
		} else {
			// If schedule component missing or disabled, ensure no active job exists for this workflow
			var job models.Jobs
			if err := db.Where("jobname = ?", wf.Wfname).First(&job).Error; err == nil {
				// We found a job with this name. If it was auto-generated by us (or matches flow name), disable it.
				if job.RecordStatus == 1 {
					if err := db.Model(&job).Update("recordstatus", 0).Error; err != nil {
						log.Printf("Failed to disable job %s: %v", wf.Wfname, err)
					} else {
						log.Printf("Disabled job %s because schedule component was removed or disabled", wf.Wfname)
					}
				}
			}
		}
	}
}

// Add this at the end of scheduler.go before the closing brace

// validateCommand performs basic validation and sanitization on the command
func validateCommand(cmd string) error {
	// Trim whitespace
	cmd = strings.TrimSpace(cmd)

	// Check if command is empty
	if cmd == "" {
		return fmt.Errorf("command is empty")
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"rm -rf /",
		"mkfs",
		"dd if=",
		":(){ :|:& };:", // Fork bomb
		"> /dev/sda",
		"chmod 777 /",
	}

	cmdLower := strings.ToLower(cmd)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmdLower, pattern) {
			return fmt.Errorf("command contains dangerous pattern: %s", pattern)
		}
	}

	// Check for command injection attempts
	if strings.Contains(cmd, "&&") || strings.Contains(cmd, "||") || strings.Contains(cmd, ";") {
		log.Printf("Warning: Command contains chaining operators: %s", cmd)
		// Allow but log - user might legitimately need these
	}

	return nil
}

// executeCommand executes an external command with validation
func executeCommand(cmdStr string) (string, error) {
	// Validate command
	if err := validateCommand(cmdStr); err != nil {
		return "", fmt.Errorf("command validation failed: %v", err)
	}

	var cmd *exec.Cmd

	// Determine OS and create appropriate command
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// Capture combined output (stdout + stderr)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return outputStr, fmt.Errorf("command failed: %v", err)
	}

	// Log output if not empty
	if len(output) > 0 {
		log.Printf("Command output: %s", outputStr)
	}

	return outputStr, nil
}
