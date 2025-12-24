package generator

/*
WORKFLOW EXECUTION ENGINE DOCUMENTATION

This file implements the core workflow execution engine for the ERP6 system.
It handles the execution of visual workflow diagrams created in the frontend designer.

=== ARCHITECTURE OVERVIEW ===

1. WORKFLOW STRUCTURE:
   - Workflows are stored as JSON in the database (Drawflow format)
   - Each workflow contains multiple nodes (components) connected by edges
   - Nodes can be: Start, End, Decision, or custom components (AI, Scraper, SendMessage, etc.)

2. EXECUTION FLOW:
   ExecuteFlow() → InternalFlow() → Component Handler → InternalFlow() (recursive)

   - ExecuteFlow: Entry point, loads workflow, initializes state, starts execution
   - InternalFlow: Executes a single node and recursively calls next connected nodes
   - Component Handlers: Registered functions that implement specific node logic

3. STATE MANAGEMENT (stored in fiber.Ctx.Locals):
   - "wfEngine": []WorkflowEngine - Tracks execution history and results of each node
   - "flowTerminated": bool - Flag to stop execution when End node is reached
   - "components": []Component - All nodes in the workflow
   - "scopedParams": map[string]interface{} - Parameters passed to nested workflows

4. DATA FLOW BETWEEN NODES:
   - Each node stores its result in wfEngine (c.Locals("wfEngine"))
   - Subsequent nodes can access previous results via ResolveParam() using $variable syntax
   - Example: If AI node outputs {"action": "scrape"}, next node can use $action

5. EXECUTION ORDER:
   - Topological sort determines initial order based on connections
   - Actual execution is recursive, following the connection graph
   - Decision nodes can branch to different paths based on conditions

6. WEBSOCKET INTEGRATION:
   - Real-time updates sent to frontend during execution
   - Events: node_start, node_complete, node_error
   - Allows live visualization of workflow execution in the designer

=== KEY FUNCTIONS ===

- ExecuteFlow(): Main entry point, orchestrates workflow execution
- InternalFlow(): Recursive function that executes nodes and follows connections
- GetWorkflowDetail(): Loads node configuration from database
- ResolveParam(): Resolves variable references ($var) to actual values
- appendStepResult(): Records node execution results
- broadcastNodeUpdate(): Sends WebSocket updates to frontend

=== COMPLEXITY NOTES ===

The current implementation uses recursive traversal which can be hard to debug.
Consider refactoring to iterative approach with explicit queue/stack for better clarity.
*/

import (
	"encoding/json"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// WorkflowDetailResult represents the configuration for a single input/output parameter of a workflow node.
// This structure is loaded from the database and contains both the component definition and user-configured values.
type WorkflowDetailResult struct {
	ComponentDetailID   int    `json:"componentdetailid"` // Unique ID for this parameter definition
	ComponentName       string `json:"componentname"`     // Name of the component (e.g., "AI", "Scraper")
	DetailType          string `json:"detailtype"`        // "input" or "output"
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

// WorkflowEngine tracks the execution state and results of a single node.
// An array of these is stored in c.Locals("wfEngine") to maintain execution history.
// This allows subsequent nodes to access results from previous nodes via ResolveParam().
type WorkflowEngine struct {
	WorkflowId    int     `json:"workflowId"`      // ID of the workflow being executed
	NodeId        int     `json:"nodeId"`          // ID of the specific node instance
	ComponentName string  `json:"componentName"`   // Type of component (e.g., "AI", "Scraper")
	DataInputNode any     `json:"input"`           // Input parameters sent to this node
	ResultNode    any     `json:"result"`          // Output/result from this node (accessible via $variable)
	Success       bool    `json:"success"`         // Whether execution succeeded
	ExecutionTime float64 `json:"executionTime"`   // Execution time in milliseconds
	Error         string  `json:"error,omitempty"` // Error message if failed
}

// Connection represents an edge between two nodes in the workflow graph.
// Defines how data flows from one node's output to another node's input.
type Connection struct {
	Node   string `json:"node"`   // Target node ID
	Output string `json:"output"` // Output port name on source node
	Input  string `json:"input"`  // Input port name on target node
}

type IO struct {
	Connections []Connection `json:"connections"`
}

// Component represents a single node in the workflow graph.
// This structure is deserialized from the Drawflow JSON stored in the database.
type Component struct {
	WorkflowId int               `json:"workflowid"` // Parent workflow ID
	ID         int               `json:"id"`         // Unique node ID within the workflow
	Name       string            `json:"name"`       // Component type (e.g., "AI", "Scraper", "Start", "End")
	Class      string            `json:"class"`      // CSS class for frontend rendering
	HTML       string            `json:"html"`       // HTML template for frontend
	Typenode   bool              `json:"typenode"`   // Whether this is a special node type
	Inputs     map[string]IO     `json:"inputs"`     // Input ports and their connections
	Outputs    map[string]IO     `json:"outputs"`    // Output ports and their connections
	PosX       float64           `json:"pos_x"`      // X position in designer canvas
	PosY       float64           `json:"pos_y"`      // Y position in designer canvas
	Data       map[string]string `json:"data"`       // Additional metadata
	IsRun      bool              `json:"isrun"`      // Runtime flag to prevent duplicate execution
}

type FlowData struct {
	Drawflow struct {
		Home struct {
			Data map[string]Component `json:"data"`
		} `json:"Home"`
	} `json:"drawflow"`
}

// GetWorkflowDetail loads the configuration for a specific node from the database.
// It retrieves both the component definition (inputs/outputs) and user-configured values.
// Returns an array of WorkflowDetailResult, one for each input/output parameter.
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

// handleStart is the component handler for the "Start" node.
// It initializes the wfEngine array to begin tracking execution.
func handleStart(c *fiber.Ctx) error {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)

	// Load workflow parameters and their values
	flowName := c.FormValue("flowname")
	params := make(map[string]interface{})

	// First, load from workflow definition (defaults)
	if flowName != "" {
		// Check if db is available
		dbInterface := c.Locals("db")
		if dbInterface != nil {
			db := dbInterface.(*gorm.DB)

			// Get workflow parameters definition
			var wfParams []models.Workflowparameter
			if err := db.
				Table("workflowparameter a").
				Select("a.parametername, a.parametervalue").
				Joins("INNER JOIN workflow b ON b.workflowid = a.workflowid").
				Where("b.wfname = ?", flowName).
				Scan(&wfParams).Error; err == nil {

				// Load parameter values from form or use default
				for _, p := range wfParams {
					val := c.FormValue(p.Parametername)
					if val == "" && p.Parametervalue != "" {
						val = p.Parametervalue
					}
					params[p.Parametername] = val
				}
			}
		}
	}

	// Then, override with parameters from parent workflow (scopedParams takes priority)
	if scopedParams, ok := c.Locals("scopedParams").(map[string]interface{}); ok {
		for k, v := range scopedParams {
			params[k] = v
		}
	}

	// Load conversation state from file (for AI assistant continuity)
	if userID, ok := c.Locals("userid").(int); ok && userID > 0 {
		conversationFile := fmt.Sprintf("./tmp/ai_conversations/%d.json", userID)
		if data, err := os.ReadFile(conversationFile); err == nil {
			var state map[string]interface{}
			if json.Unmarshal(data, &state) == nil {
				if convState, ok := state["conversation_state"].(string); ok && convState != "" {
					params["conversation_state"] = convState
					fmt.Printf("[Start Node] Loaded conversation state for user %d\n", userID)
				}
			}
		}
	}

	// Debug: Log what parameters we're setting
	fmt.Printf("[Start Node] Parameters loaded: %+v\n", params)

	// Make parameters available to subsequent nodes
	wfEngine = append(wfEngine, WorkflowEngine{
		DataInputNode: "",
		ResultNode:    params,
	})
	c.Locals("wfEngine", wfEngine)
	return nil
}

// GetSearchText retrieves parameter values from various sources (POST, GET, scopedParams).
// Used primarily for search/filter operations with special handling for date/time fields.
// Converts string values to LIKE patterns (%value%) for database queries.
func GetSearchText(c *fiber.Ctx, paramTypes []string, param, defVal, dataType string) string {
	s := defVal

	for _, t := range paramTypes {
		switch strings.ToUpper(t) {
		case "POST":
			// Check scoped params first (simulating POST/local scope)
			if scopedParams, ok := c.Locals("scopedParams").(map[string]interface{}); ok {
				if val, exists := scopedParams[param]; exists {
					s = fmt.Sprintf("%v", val)
					break
				}
			}

			if val := c.FormValue(param); val != "" {
				s = val
			}
		case "GET":
			// Check scoped params for GET too? Usually safer to allow it to override both.
			if scopedParams, ok := c.Locals("scopedParams").(map[string]interface{}); ok {
				if val, exists := scopedParams[param]; exists {
					s = fmt.Sprintf("%v", val)
					break
				}
			}

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

// InternalFlow executes a single workflow node and recursively processes connected nodes.
// This is the core execution function that:
// 1. Checks if node already executed (prevents loops)
// 2. Loads node configuration from database
// 3. Executes the component handler
// 4. Records results in wfEngine
// 5. Broadcasts WebSocket updates
// 6. Recursively calls itself for connected nodes
//
// COMPLEXITY WARNING: This recursive approach can be hard to debug.
// Consider refactoring to iterative approach with explicit queue.
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

	// Broadcast node start via WebSocket
	broadcastNodeUpdate(c, "node_start", workflowId, nodeId, component.Name, nil, nil, 0, "")

	// Small delay to allow frontend to process the event and update UI
	// This ensures fast-executing components (like Transform) show the "running" state
	time.Sleep(50 * time.Millisecond)

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

	// Retrieve extras from Fiber Locals if available (persistence across nodes)
	extras, _ := c.Locals("wfExtras").(map[string]interface{})
	if extras == nil {
		extras = make(map[string]interface{})
	}

	// Create context for the component
	ctx := &WorkflowContext{
		FiberCtx:         c,
		DB:               db,
		Params:           workflowDetailResult,
		Search:           search,
		CurrentComponent: component,
		Extras:           extras,
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
		// Save extras back to locals to persist changes
		c.Locals("wfExtras", ctx.Extras)
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

		// Broadcast node completion for external plugins
		broadcastNodeUpdate(c, "node_complete", workflowId, nodeId, component.Name, inputParams, stepResult, execTime, "")
	} else {
		// Component didn't append, add our own tracking
		appendStepResult(c, workflowId, nodeId, component.Name, inputParams, stepResult, true, execTime, "")
	}

	// Check if navigation should be skipped (handled manually by component)
	if skip, ok := c.Locals("skipNavigation").(bool); ok && skip {
		// Reset flag to avoid affecting parent calls if context is reused strangely
		// (though in recursion, we return immediately, so this is mostly for safety)
		c.Locals("skipNavigation", false)
		return nil
	}

	// Handle Decision and Loop flow
	if strings.EqualFold(component.Name, "Decision") || strings.EqualFold(component.Name, "auth") || strings.EqualFold(component.Name, "For Each") {
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

// appendStepResult adds a node execution record to the wfEngine history.
// This allows subsequent nodes to access results via ResolveParam($variable).
// Also broadcasts WebSocket updates to the frontend for live visualization.
func appendStepResult(c *fiber.Ctx, workflowId int, nodeId int, componentName string, input any, result any, success bool, execTime float64, errMsg string) {
	// Persistent Debug Logging (Production Safe)
	fmt.Printf("[Workflow Debug] Node: %s (ID: %d) | Input: %+v | Result: %+v | Success: %v | Error: %s\n", 
		componentName, nodeId, input, result, success, errMsg)

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

	// Broadcast node completion via WebSocket
	eventType := "node_complete"
	if !success {
		eventType = "node_error"
	}
	broadcastNodeUpdate(c, eventType, workflowId, nodeId, componentName, input, result, execTime, errMsg)
}

// broadcastNodeUpdate sends workflow node updates via WebSocket
// broadcastNodeUpdate sends real-time workflow execution updates via WebSocket.
// Events: node_start, node_complete, node_error
// Allows the frontend designer to visualize execution progress in real-time.
func broadcastNodeUpdate(c *fiber.Ctx, eventType string, workflowId int, nodeId int, componentName string, input any, result any, execTime float64, errMsg string) {
	if ws.GlobalHub == nil {
		return
	}

	// Get user ID from context
	userID, ok := c.Locals("userid").(int)
	if !ok {
		return
	}

	payload := map[string]interface{}{
		"type":          "workflow_test",
		"event":         eventType,
		"workflowId":    workflowId,
		"nodeId":        nodeId,
		"componentName": componentName,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	if eventType != "node_start" {
		payload["input"] = input
		payload["result"] = result
		payload["executionTime"] = execTime
		payload["success"] = errMsg == ""
		if errMsg != "" {
			payload["error"] = errMsg
		}
	}

	message, err := json.Marshal(payload)
	if err == nil {
		ws.GlobalHub.SendToUser(userID, message)
	}
}

// ExecuteFlow is the main entry point for workflow execution.
// It orchestrates the entire workflow execution process:
// 1. Loads workflow definition from database
// 2. Parses Drawflow JSON structure
// 3. Performs topological sort to determine execution order
// 4. Initializes execution state (wfEngine, components, etc.)
// 5. Starts execution from the first node (usually "Start")
// 6. Manages database transaction (commit/rollback)
//
// Parameters:
//   - flowName: Name of the workflow to execute
//   - search: If true, runs in read-only mode (no transaction)
//   - params: Optional parameters to pass to the workflow (stored in scopedParams)
//
// The actual node-by-node execution is handled by InternalFlow() which is called recursively.
func ExecuteFlow(c *fiber.Ctx, db *gorm.DB, flowName string, search bool, params map[string]interface{}) error {
	// Stack management for scopedParams
	originalScopedParams := c.Locals("scopedParams")
	if params != nil {
		c.Locals("scopedParams", params)
	}
	defer c.Locals("scopedParams", originalScopedParams)

	// Store db in locals for component access
	c.Locals("db", db)

	var components = []Component{}
	var wfEngine = []WorkflowEngine{}
	var flowTerminated = false

	// Get workflow parameters
	var wfParams []models.Workflowparameter
	if err := db.
		Table("workflowparameter a").
		Select("a.wfparameterid, b.wfname, a.parametername").
		Joins("INNER JOIN workflow b ON b.workflowid = a.workflowid").
		Where("b.wfname = ?", flowName).
		Scan(&wfParams).Error; err != nil {
		return err
	}

	// Initialize default parameters
	postData := make(map[string]interface{})
	for _, p := range wfParams {
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
		// Check if we're in a nested workflow (called from Workflow component)
		isNested, _ := c.Locals("nestedWorkflow").(bool)

		if isNested {
			// Nested workflow - reuse the existing transaction
			tx = db
			fmt.Printf("[ExecuteFlow] Reusing existing transaction for nested workflow\n")
		} else {
			// Top-level workflow - start new transaction
			tx = db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					fmt.Printf("Transaction panic: %v\n", r)
				}
			}()
			fmt.Printf("[ExecuteFlow] Started new transaction\n")
		}
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

	// Commit transaction (only for top-level workflows)
	if !search {
		isNested, _ := c.Locals("nestedWorkflow").(bool)
		if !isNested {
			// Top-level workflow - commit the transaction
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("commit failed: %w", err)
			}
			fmt.Printf("[ExecuteFlow] Transaction committed\n")
		} else {
			// Nested workflow - don't commit, let parent handle it
			fmt.Printf("[ExecuteFlow] Skipping commit for nested workflow\n")
		}
	}

	return nil
}

// ResolveParam resolves a parameter value, supporting variable substitution (e.g., $param).
// This is the key function that enables data flow between workflow nodes.
//
// Resolution order:
// 1. If value starts with $, treat as variable reference
// 2. Check POST form values (c.FormValue)
// 3. Check GET query parameters (c.Query)
// 4. Check previous node results in wfEngine (searches backwards for latest match)
//   - Checks top-level keys in ResultNode
//   - Checks nested "data" object in ResultNode
//
// 5. If not found, return original value
//
// Example: If AI node outputs {"action": "scrape", "url": "example.com"},
// subsequent nodes can use $action and $url to access these values.
func ResolveParam(c *fiber.Ctx, val string) string {
	// Handle embedded variables like: https://example.com?q=$city_name
	// Supports recursive resolution (up to 5 levels)
	if strings.Contains(val, "$") {
		result := val
		for i := 0; i < 5; i++ { // limit recursion
			if !strings.Contains(result, "$") {
				break
			}

			// Find all $variable patterns
			re := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)
			matches := re.FindAllStringSubmatch(result, -1)

			changed := false
			for _, match := range matches {
				if len(match) >= 2 {
					varName := match[1]
					varValue := resolveVariable(c, varName)

					// Only replace if value is different (prevent infinite loops with unresolved vars)
					if varValue != "$"+varName {
						result = strings.ReplaceAll(result, "$"+varName, varValue)
						changed = true
					}
				}
			}

			if !changed {
				break
			}
		}
		return result
	}

	return val
}

// resolveVariable looks up a variable value from various sources
func resolveVariable(c *fiber.Ctx, key string) string {
	// 1. Check POST
	if v := c.FormValue(key); v != "" {
		return v
	}

	// 2. Check GET
	if v := c.Query(key); v != "" {
		return v
	}

	// 2.5 Check Scoped Params (from Iterator/Schedule)
	if scopedParams, ok := c.Locals("scopedParams").(map[string]interface{}); ok {
		if v, exists := scopedParams[key]; exists {
			return fmt.Sprintf("%v", v)
		}
	}

	// 2.6 Check Extras (from For Each / Manual Set)
	if extras, ok := c.Locals("wfExtras").(map[string]interface{}); ok {
		if v, exists := extras[key]; exists {
			return fmt.Sprintf("%v", v)
		}
	}

	// 3. Check node results
	wfEngine := c.Locals("wfEngine")
	if wfEngine != nil {
		if engines, ok := wfEngine.([]WorkflowEngine); ok {
			// Iterate backwards to find latest match
			for i := len(engines) - 1; i >= 0; i-- {
				e := engines[i]
				if resMap, ok := e.ResultNode.(map[string]interface{}); ok {
					// Check top level
					if v, exists := resMap[key]; exists {
						return fmt.Sprintf("%v", v)
					}
					// Check inside "data"
					if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
						if v, exists := dataMap[key]; exists {
							return fmt.Sprintf("%v", v)
						}
					}
				}
			}
		}
	}

	// 4. Check Locals (auth context etc)
	if v := c.Locals(key); v != nil {
		return fmt.Sprintf("%v", v)
	}

	// Return original key if not found
	return "$" + key
}
