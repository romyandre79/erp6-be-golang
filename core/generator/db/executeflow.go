package generator

import (
	"encoding/json"
	"erp6-be-golang/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WorkflowDetailResult struct {
	ComponentDetailID   int    `json:"componentdetailid"`
	ComponentName       string `json:"componentname"`
	DetailType          string `json:"detailtype"`
	Label               string `json:"label"`
	InputType           string `json:"inputtype"`
	InputName           string `json:"inputname"`
	InputDesc           string `json:"inputdesc"`
	DataSourceType      string `json:"datasourcetype"`
	DataSource          string `json:"datasource"`
	DataSourceIDField   string `json:"datasourceidfield"`
	DataSourceNameField string `json:"datasourcenamefield"`
	CompValue           string `json:"compvalue"`
	WfDetailID          int    `json:"wfdetailid"`
}

type WorkflowEngine struct {
	WorkflowId    int     `json:"workflowId"`
	NodeId        int     `json:"nodeId"`
	ComponentName string  `json:"componentName"`
	DataInputNode any     `json:"input"`
	ResultNode    any     `json:"result"`
	Success       bool    `json:"success"`
	ExecutionTime float64 `json:"executionTime"`
	Error         string  `json:"error,omitempty"`
}

type Connection struct {
	Node   string `json:"node"`
	Output string `json:"output"`
	Input  string `json:"input"`
}

type IO struct {
	Connections []Connection `json:"connections"`
}

type Component struct {
	WorkflowId int               `json:"workflowid"`
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Class      string            `json:"class"`
	HTML       string            `json:"html"`
	Typenode   bool              `json:"typenode"`
	Inputs     map[string]IO     `json:"inputs"`
	Outputs    map[string]IO     `json:"outputs"`
	PosX       float64           `json:"pos_x"`
	PosY       float64           `json:"pos_y"`
	Data       map[string]string `json:"data"`
	IsRun      bool              `json:"isrun"`
}

type FlowData struct {
	Drawflow struct {
		Home struct {
			Data map[string]Component `json:"data"`
		} `json:"Home"`
	} `json:"drawflow"`
}

func GetWorkflowDetail(db *gorm.DB, componentName string, workflowID int, nodeID int) ([]WorkflowDetailResult, error) {
	var results []WorkflowDetailResult

	var component models.Component
	err := db.
		Preload("Details.WorkflowDetails", "workflowid = ? AND nodeid = ?", workflowID, nodeID).
		Where("componentname = ?", componentName).
		First(&component).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []WorkflowDetailResult{}, nil
		}
		return nil, fmt.Errorf("failed to load component: %w", err)
	}

	for _, d := range component.Details {
		var compValue string
		var wfDetailID int

		if len(d.WorkflowDetails) > 0 {
			compValue = d.WorkflowDetails[0].Componentvalue
			wfDetailID = d.WorkflowDetails[0].Workflowdetailid
		}

		results = append(results, WorkflowDetailResult{
			ComponentDetailID:   int(d.Componentdetailid),
			ComponentName:       component.Componentname,
			DetailType:          d.Detailtype,
			Label:               d.Lable,
			InputType:           d.Inputtype,
			InputName:           d.Inputname,
			InputDesc:           d.Inputdesc,
			DataSourceType:      d.Datasourcetype,
			DataSource:          d.Datasource,
			DataSourceIDField:   d.Datasourceidfield,
			DataSourceNameField: d.Datasourcenamefield,
			CompValue:           compValue,
			WfDetailID:          wfDetailID,
		})
	}

	return results, nil
}

func handleStart(c *fiber.Ctx) error {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: ""})
	c.Locals("wfEngine", wfEngine)
	return nil
}

func GetSearchText(c *fiber.Ctx, paramTypes []string, param, defVal, dataType string) string {
	s := defVal

	for _, t := range paramTypes {
		switch strings.ToUpper(t) {
		case "POST":
			if val := c.FormValue(param); val != "" {
				s = val
			}
		case "GET":
			if val := c.Query(param); val != "" {
				s = val
			}
		case "Q":
			var val string
			if val = c.Query("q"); val != "" {
				s = val
			} else if val = c.FormValue("q"); val != "" {
				s = val
			}
			fmt.Printf("%s", val)
		}

		// Format date, datetime, time
		if strings.Contains(strings.ToLower(param), "date") && !strings.Contains(strings.ToLower(param), "datetime") {
			if s != "" {
				t, err := time.Parse("02-01-2006", s)
				if err == nil {
					s = t.Format("2006-01-02")
				}
			}
		} else if strings.Contains(strings.ToLower(param), "datetime") {
			if s != "" {
				t, err := time.Parse("02-01-2006 15:04:05", s)
				if err == nil {
					s = t.Format("2006-01-02 15:04:05")
				}
			}
		} else if strings.Contains(strings.ToLower(param), "time") {
			if s != "" {
				t, err := time.Parse("15:04:05", s)
				if err == nil {
					s = t.Format("15:04:05")
				}
			}
		}
	}

	// If datatype is string → convert to LIKE pattern: %word%
	if strings.ToLower(dataType) == "string" {
		s = "%" + strings.ReplaceAll(strings.TrimSpace(s), " ", "%") + "%"
	}

	return s
}

func getUserObjectValues(db *gorm.DB, username string, menuobject string) (string, error) {
	var ids []string
	err := db.Table("groupmenuauth AS a").
		Select("DISTINCT a.menuvalueid").
		Joins("INNER JOIN menuauth b ON b.menuauthid = a.menuauthid").
		Joins("INNER JOIN usergroup c ON c.groupaccessid = a.groupaccessid").
		Joins("INNER JOIN useraccess d ON d.useraccessid = c.useraccessid").
		Where("b.menuobject = ? AND d.username = ?", menuobject, username).
		Pluck("a.menuvalueid", &ids).Error

	if err != nil {
		return "", err
	}

	return strings.Join(ids, ","), nil
}

func getDataByCompany(db *gorm.DB, username string, dataType string) (string, error) {
	cid := ""
	retSql := make(map[string]interface{})
	dataObjectValue, _ := getUserObjectValues(db, username, "company")
	dataSplit := strings.Split(dataObjectValue, ",")
	for _, v := range dataSplit {
		sqlStat := "select a.addressbookid from addressbook a join addresscompany b on b.addressbookid = a.addressbookid where " + dataType + " = 1 and companyid = " + v
		if err := db.Raw(sqlStat).Scan(&retSql).Error; err != nil {
			return "", err
		}
		for _, v := range retSql {
			if cid == "" {
				cid = fmt.Sprintf("%+v", v)
			} else {
				cid += "," + fmt.Sprintf("%+v", v)
			}
		}
	}
	return cid, nil
}

// Component registration - only native Start and End
func init() {
	RegisterComponent("start", func(ctx *WorkflowContext) error {
		return handleStart(ctx.FiberCtx)
	})
	RegisterComponent("end", func(ctx *WorkflowContext) error {
		return nil
	})
}

func InternalFlow(c *fiber.Ctx, component Component, workflowId int, nodeId int, db *gorm.DB, search bool) error {
	var flowTerminated = c.Locals("flowTerminated").(bool)
	var components = c.Locals("components").([]Component)
	if flowTerminated {
		return nil
	}

	if component.IsRun {
		return nil
	}

	component.IsRun = true
	fmt.Printf("Running workflowid: %d component: %s (ID: %d)\n", workflowId, component.Name, component.ID)

	// Start timing
	startTime := time.Now()

	workflowDetailResult, err := GetWorkflowDetail(db, component.Name, workflowId, nodeId)
	if err != nil {
		// Record failed step
		appendStepResult(c, workflowId, nodeId, component.Name, nil, nil, false, 0, err.Error())
		return err
	}

	// Prepare input params for tracking
	inputParams := make(map[string]string)
	for _, p := range workflowDetailResult {
		if p.CompValue != "" {
			inputParams[p.InputName] = p.CompValue
		}
	}

	// Create context for the component
	ctx := &WorkflowContext{
		FiberCtx: c,
		DB:       db,
		Params:   workflowDetailResult,
		Search:   search,
	}

	// Handle special cases for context population
	if strings.ToLower(component.Name) == "importdata" {
		file, _ := c.FormFile("file-modules")
		ctx.FileHeader = file
	}

	// Execute component
	if strings.EqualFold(component.Name, "End") {
		flowTerminated = true
		c.Locals("flowTerminated", true)
		execTime := float64(time.Since(startTime).Milliseconds())
		appendStepResult(c, workflowId, nodeId, component.Name, inputParams, "Flow ended", true, execTime, "")
		return nil
	}

	handler, exists := GetComponent(component.Name)
	if exists {
		if err := handler.Execute(ctx); err != nil {
			execTime := float64(time.Since(startTime).Milliseconds())
			appendStepResult(c, workflowId, nodeId, component.Name, inputParams, nil, false, execTime, err.Error())
			return err
		}
	} else {
		execTime := float64(time.Since(startTime).Milliseconds())
		appendStepResult(c, workflowId, nodeId, component.Name, inputParams, nil, false, execTime, fmt.Sprintf("unknown component: %s", component.Name))
		return fmt.Errorf("unknown component: %s", component.Name)
	}

	// Record execution time and get the last result from wfEngine
	execTime := float64(time.Since(startTime).Milliseconds())

	// Get any result that might have been set by the component
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	var stepResult any = "OK"
	if len(wfEngine) > 0 {
		lastResult := wfEngine[len(wfEngine)-1]
		if lastResult.ResultNode != nil {
			stepResult = lastResult.ResultNode
		}
	}

	// Update the last step with component info if it was just added, or add new one
	if len(wfEngine) > 0 && wfEngine[len(wfEngine)-1].ComponentName == "" {
		// Component already appended result, update with metadata
		wfEngine[len(wfEngine)-1].WorkflowId = workflowId
		wfEngine[len(wfEngine)-1].NodeId = nodeId
		wfEngine[len(wfEngine)-1].ComponentName = component.Name
		wfEngine[len(wfEngine)-1].DataInputNode = inputParams
		wfEngine[len(wfEngine)-1].Success = true
		wfEngine[len(wfEngine)-1].ExecutionTime = execTime
		c.Locals("wfEngine", wfEngine)
	} else {
		// Component didn't append, add our own tracking
		appendStepResult(c, workflowId, nodeId, component.Name, inputParams, stepResult, true, execTime, "")
	}

	// Handle Decision flow
	if strings.EqualFold(component.Name, "Decision") {
		var outputs IO
		if ctx.DecisionResult {
			outputs = component.Outputs["output_1"]
		} else {
			outputs = component.Outputs["output_2"]
		}

		for _, conn := range outputs.Connections {
			nextNodeId, _ := strconv.Atoi(conn.Node)
			for _, nextComp := range components {
				if nextComp.ID == nextNodeId && nextComp.WorkflowId == workflowId {
					return InternalFlow(c, nextComp, workflowId, nextNodeId, db, search)
				}
			}
		}
		return nil
	} else {
		// Standard flow
		for _, output := range component.Outputs {
			for _, conn := range output.Connections {
				nextNodeId, _ := strconv.Atoi(conn.Node)
				for _, nextComp := range components {
					if nextComp.ID == nextNodeId && nextComp.WorkflowId == workflowId {
						err := InternalFlow(c, nextComp, workflowId, nextNodeId, db, search)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// appendStepResult adds a step result to the workflow engine
func appendStepResult(c *fiber.Ctx, workflowId int, nodeId int, componentName string, input any, result any, success bool, execTime float64, errMsg string) {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{
		WorkflowId:    workflowId,
		NodeId:        nodeId,
		ComponentName: componentName,
		DataInputNode: input,
		ResultNode:    result,
		Success:       success,
		ExecutionTime: execTime,
		Error:         errMsg,
	})
	c.Locals("wfEngine", wfEngine)
}

func ExecuteFlow(c *fiber.Ctx, db *gorm.DB, flowName string, search bool) error {
	var components = []Component{}
	var wfEngine = []WorkflowEngine{}
	var flowTerminated = false

	// Get workflow parameters
	var params []models.Workflowparameter
	if err := db.
		Table("workflowparameter a").
		Select("a.wfparameterid, b.wfname, a.parametername").
		Joins("INNER JOIN workflow b ON b.workflowid = a.workflowid").
		Where("b.wfname = ?", flowName).
		Scan(&params).Error; err != nil {
		return err
	}

	// Initialize default parameters
	postData := make(map[string]interface{})
	for _, p := range params {
		postData[p.Parametername] = nil
	}

	// Get workflow definition
	var wf models.Workflow
	if err := db.
		Where("wfname = ?", flowName).
		First(&wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("workflow %s does not exist", flowName)
		}
		return errors.New("INVALID_WORKFLOW")
	}

	// Decode JSON flow
	var flow FlowData
	if err := json.Unmarshal([]byte(wf.Flow), &flow); err != nil {
		return fmt.Errorf("failed to decode flow JSON: %w", err)
	}

	// Convert map to slice
	for _, comp := range flow.Drawflow.Home.Data {
		comp.WorkflowId = wf.Workflowid
		components = append(components, comp)
	}

	// Build connection graph
	graph := make(map[int][]int)
	inDegree := make(map[int]int)
	for _, comp := range components {
		for _, out := range comp.Outputs {
			for _, conn := range out.Connections {
				src := comp.ID
				dst, _ := strconv.Atoi(conn.Node)
				graph[src] = append(graph[src], dst)
				inDegree[dst]++
				if _, ok := inDegree[src]; !ok {
					inDegree[src] = 0
				}
			}
		}
	}

	// Find start nodes (inDegree == 0)
	queue := []int{}
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	// Topological sort
	ordered := []Component{}
	visited := make(map[int]bool)

	for len(queue) > 0 {
		currID := queue[0]
		queue = queue[1:]

		if visited[currID] {
			continue
		}
		visited[currID] = true

		for _, comp := range components {
			if comp.ID == currID {
				comp.WorkflowId = wf.Workflowid
				ordered = append(ordered, comp)
				break
			}
		}

		for _, next := range graph[currID] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	// Add disconnected nodes
	if len(ordered) < len(components) {
		for _, comp := range components {
			if !visited[comp.ID] {
				ordered = append(ordered, comp)
			}
		}
	}

	// Setup database transaction
	var tx *gorm.DB
	if !search {
		tx = db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				fmt.Printf("Transaction panic: %v\n", r)
			}
		}()
	} else {
		tx = db
	}

	c.Locals("components", ordered)
	c.Locals("wfEngine", wfEngine)
	c.Locals("flowTerminated", flowTerminated)

	// Execute workflow starting from first component
	if err := InternalFlow(c, ordered[0], int(wf.Workflowid), ordered[0].ID, tx, search); err != nil {
		if !search {
			tx.Rollback()
		}
		return fmt.Errorf("internal flow failed on component %d (%s): %w", ordered[0].ID, ordered[0].Name, err)
	}

	// Commit transaction
	if !search {
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
	}

	return nil
}
