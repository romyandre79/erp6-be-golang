package scheduler

import (
	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/models"
	"log"
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
		var wfEngine []generator.WorkflowEngine // Using the package name 'generator' alias for 'erp6-be-golang/core/generator/db'?
		// Wait, the package declaration in executeflow.go is "package generator".
		// So import "erp6-be-golang/core/generator/db" will likely be aliased as "db" or "generator"?
		// usually last folder name unless specified. folder is "db", file package is "generator".
		// I should alias it to avoid confusion or if it's main.

		ctx.Locals("wfEngine", wfEngine)
		ctx.Locals("flowTerminated", false)

		if err := generator.ExecuteFlow(ctx, gormDB, job.Flow, false); err != nil {
			log.Printf("Error executing flow %s: %v", job.Flow, err)
		} else {
			log.Printf("Flow %s executed successfully", job.Flow)
		}
	}
}
