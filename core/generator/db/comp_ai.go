package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// -- Structs converted from component-ai/main.go --

// AIConversationState represents the current state of a conversation
type AIConversationState struct {
	EntityType          string            `json:"entity_type"`          // module, menu, table, workflow
	CurrentStep         int               `json:"current_step"`         // Current question index
	CollectedData       map[string]string `json:"collected_data"`       // Data collected so far
	IsComplete          bool              `json:"is_complete"`          // Whether conversation is complete
	WaitingConfirmation bool              `json:"waiting_confirmation"` // Waiting for user to type 'execute'
}

// AIQuestionFlow defines the questions for each entity type and query templates
type AIQuestionFlow struct {
	Questions      []AIQuestion      `json:"questions"`
	Queries        map[string]string `json:"queries"`
	ListQueries    map[string]string `json:"list_queries"`
	Params         []string          `json:"params"`
	Defaults       map[string]string `json:"defaults"`
	SuccessMessage string            `json:"success_message"`
	MetaAction     string            `json:"meta_action"`
	Description    string            `json:"description"`
	Triggers       []string          `json:"triggers"`
}

type AIQuestion struct {
	Key          string   `json:"key"`           // Field name to store answer
	Text         string   `json:"text"`          // Question to ask user
	Validation   string   `json:"validation"`    // Validation type: required, optional, number, etc.
	Options      []string `json:"options"`       // For select-type questions
	OptionsQuery string   `json:"options_query"` // SQL query to fetch options dynamically
	Description  string   `json:"description"`   // Help text
}

// EntityInfo holds information about an available command entity
type AIEntityInfo struct {
	Name        string
	Description string
	Triggers    string
}

func init() {
	RegisterComponent("AIAssistant", func(ctx *WorkflowContext) error {
		return handleAI(ctx)
	})
}

func handleAI(ctx *WorkflowContext) error {
	var (
		command           string
		userID            string
		conversationState string
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		// Fallback
		if val == "" {
			val = ResolveParam(ctx.FiberCtx, p.InputName)
			if strings.HasPrefix(val, "$") {
				val = ""
			}
		}

		switch strings.ToLower(p.InputName) {
		case "command", "message", "text", "input":
			command = val
		case "user_id":
			userID = val
		case "conversation_state":
			conversationState = val
		}
	}

	// Default user_id from context if not provided
	if userID == "" {
		if ctx.FiberCtx != nil {
			if uid, ok := ctx.FiberCtx.Locals("userid").(int); ok {
				userID = fmt.Sprintf("%d", uid)
			}
		}
	}

	// DB Driver from context/config. In erp6-be-golang it's accessible via ctx.DB.Dialector.Name() usually?
	// Or we can get it from Config.
	dbDriver := ctx.DB.Dialector.Name() // e.g. "mysql", "postgres"

	result, err := processAI(command, conversationState, dbDriver, userID, ctx.DB)
	if err != nil {
		return err
	}

	// Append result to workflow engine
	wm := WorkflowEngine{
		ResultNode: result,
	}
	
	if ctx.FiberCtx != nil {
		wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
	} else {
		// For non-HTTP contexts (e.g., WhatsApp), store in Extras
		wfEngine, _ := ctx.Extras["wfEngine"].([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.Extras["wfEngine"] = wfEngine
	}

	return nil
}

// ExecuteAIResult executes the workflow or SQL from an AI result
// Can be called from both HTTP (Fiber) and non-HTTP (WhatsApp) contexts
func ExecuteAIResult(aiResult map[string]interface{}, db *gorm.DB, ctx *WorkflowContext) error {
	exec, _ := aiResult["execute"].(string)
	if exec != "true" {
		return nil // Nothing to execute
	}
	
	metaAction, _ := aiResult["meta_action"].(string)
	
	// SQL Execution
	if metaAction == "insert" || metaAction == "update" || metaAction == "delete" {
		query, _ := aiResult["query"].(string)
		paramsStr, _ := aiResult["parameters"].(string)
		
		var params []interface{}
		if paramsStr != "" && paramsStr != "[]" {
			json.Unmarshal([]byte(paramsStr), &params)
		}
		
		if err := db.Exec(query, params...).Error; err != nil {
			return err
		}
		
		// Store success in context
		if ctx != nil {
			ctx.Extras["sql_executed"] = true
		}
		return nil
	}
	
	// Check if AI returned a workflow name to execute
	if workflowName, ok := aiResult["workflow_name"].(string); ok && workflowName != "" {
		fmt.Printf("[ExecuteAIResult] Executing workflow: %s\n", workflowName)
		
		// For WhatsApp context, we need to create a mock Fiber context
		// or execute the workflow differently
		if ctx.FiberCtx == nil {
			// WhatsApp context - execute workflow components manually
			// This is complex, so for now we'll fall back to single component execution
			fmt.Printf("[ExecuteAIResult] WhatsApp workflow execution not yet fully implemented\n")
			// Fall through to single component execution
		} else {
			// HTTP context - use ExecuteFlow
			params := make(map[string]interface{})
			for key, value := range aiResult {
				if key != "message" && key != "execute" && key != "meta_action" && 
				   key != "conversation_state" && key != "user_id" && key != "query" && 
				   key != "parameters" && key != "workflow_name" {
					params[key] = value
				}
			}
			return ExecuteFlow(ctx.FiberCtx, db, workflowName, false, params)
		}
	}
	
	// Workflow Execution - try to map to known workflows first
	if ctx == nil {
		return fmt.Errorf("workflow execution requires context")
	}
	
	// Check if we can map this to a known workflow based on patterns
	// This allows AI component execution to use full workflows
	if action, ok := aiResult["action"].(string); ok {
		if action == "extract_one_data" || action == "extract_data" {
			// Check URL to determine which workflow
			if url, ok := aiResult["url"].(string); ok {
				workflowName := ""
				if strings.Contains(url, "bi.go.id") && strings.Contains(url, "kurs") {
					workflowName = "kurs bi"
				}
				// Add more URL patterns here as needed
				
				if workflowName != "" {
					fmt.Printf("[ExecuteAIResult] Mapped scraping request to workflow: %s\n", workflowName)
					// Try to execute the full workflow
					// For WhatsApp, we need to manually execute workflow components
					// since we don't have a Fiber context
					// For now, fall through to single component execution
					// TODO: Implement full workflow execution for WhatsApp
				}
			}
		}
	}
	
	// Build params for workflow component (skip metadata)
	ctx.Params = []WorkflowDetailResult{}
	for key, value := range aiResult {
		if key != "message" && key != "execute" && key != "meta_action" && 
		   key != "conversation_state" && key != "user_id" && key != "query" && key != "parameters" {
			ctx.Params = append(ctx.Params, WorkflowDetailResult{
				InputName: key,
				CompValue: fmt.Sprintf("%v", value),
			})
		}
	}
	
	// Determine component name
	componentName := ""
	if action, ok := aiResult["action"].(string); ok {
		componentName = action
		// Map known actions
		if action == "extract_one_data" || action == "extract_data" {
			componentName = "Web Scraper"
		}
	}
	
	if componentName == "" {
		return fmt.Errorf("no component name in AI result")
	}
	
	// Execute component
	handler, ok := GetComponent(componentName)
	if !ok {
		return fmt.Errorf("component '%s' not found", componentName)
	}
	
	return handler.Execute(ctx)
}

func processAI(command, stateJSON, dbDriver, userID string, db *gorm.DB) (map[string]interface{}, error) {
	var state AIConversationState

	// Parse existing state or create new one
	// Parse existing state or create new one
	if stateJSON != "" && stateJSON != "{}" {
		if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
			state = AIConversationState{CollectedData: make(map[string]string)}
		} else {
			// Ensure map is initialized if json had null
			if state.CollectedData == nil {
				state.CollectedData = make(map[string]string)
			}
		}
	} else {
		state = AIConversationState{CollectedData: make(map[string]string)}
	}

	lowerCmd := strings.ToLower(command)

	// Check if this is a help command
	if strings.HasPrefix(lowerCmd, "help") {
		availableEntities := getAvailableEntities(db)

		var helpMsg strings.Builder
		helpMsg.WriteString("Available commands:\n")
		for _, entity := range availableEntities {
			helpMsg.WriteString(fmt.Sprintf("• %s - %s\n", entity.Triggers, entity.Description))
		}
		helpMsg.WriteString("• help - Show this help\n")
		helpMsg.WriteString("• get data [table] - Show data from a table (e.g., 'get data customer')\n")

		return map[string]interface{}{
			"message": helpMsg.String(),
		}, nil
	}

	// WA Registration Command
	if strings.HasPrefix(lowerCmd, "register wa api") {
		qr, err := GetLoginQR()
		if err != nil {
			return map[string]interface{}{
				"message": fmt.Sprintf("Error initializing WhatsApp: %v", err),
			}, nil
		}
		if qr == "Already logged in" {
			return map[string]interface{}{
				"message": "✅ WhatsApp is already connected!",
			}, nil
		}
		return map[string]interface{}{
			"message": fmt.Sprintf("Please scan this QR Code to connect:\n\n%s", qr),
			"qr":      qr, // Client can render this if needed
		}, nil
	}

	// Check for "get/show/list [table]" command (Data Retrieval)
	if (state.EntityType == "") && db != nil {
		if strings.HasPrefix(lowerCmd, "get data ") || strings.HasPrefix(lowerCmd, "ambil data ") || strings.HasPrefix(lowerCmd, "show data ") || strings.HasPrefix(lowerCmd, "list data ") {
			parts := strings.Fields(lowerCmd)
			if len(parts) >= 3 {
				tableName := parts[2]
				result, err := handleGetData(db, tableName)
				if err != nil {
					return map[string]interface{}{"error": err.Error()}, nil
				}
				
				response := map[string]interface{}{"result": result}
				
				// Check if this is a WhatsApp request (userID will be set for WA)
				// For WhatsApp, generate and send Excel file
				fmt.Printf("[AI Debug] userID: '%s', checking for Excel generation\n", userID)
				if userID != "" && userID != "0" {
					if list, ok := result.([]map[string]interface{}); ok && len(list) > 0 {
						fmt.Printf("[AI Debug] Generating Excel for %d rows\n", len(list))
						// Generate Excel file
						filePath, err := generateExcelFromData(tableName, list)
						if err == nil {
							fmt.Printf("[AI Debug] Excel generated: %s\n", filePath)
							// Send via WhatsApp
							// We need the phone number - it should be passed in context
							// For now, store file path in response for WhatsApp handler to send
							response["excel_file"] = filePath
							response["message"] = fmt.Sprintf("📊 Data from '%s' table (%d rows)\nSending as Excel file...", tableName, len(list))
						} else {
							fmt.Printf("[AI Debug] Excel generation failed: %v\n", err)
							// Fallback to text if Excel generation fails
							response["message"] = formatDataAsTable(list)
						}
					} else {
						response["message"] = fmt.Sprintf("No data found in table '%s'", tableName)
					}
				} else {
					fmt.Printf("[AI Debug] Not WhatsApp context (userID='%s'), using text format\n", userID)
					// Web interface - use text format
					if list, ok := result.([]map[string]interface{}); ok {
						response["message"] = formatDataAsTable(list)
					}
				}
				
				return response, nil
			}
		}
	}
	// Continue existing conversation
	if state.EntityType != "" {
		return runConversationStep(command, state, dbDriver, userID, "", "", db)
	}

	// New Conversation Logic
	availableEntities := getAvailableEntities(db)
	matchedEntity := ""
	initialArg := ""
	isListCommand := false

	if strings.HasPrefix(lowerCmd, "list ") {
		for _, entity := range availableEntities {
			if strings.HasPrefix(lowerCmd, "list "+entity.Name) || strings.HasPrefix(lowerCmd, "list "+entity.Name+"s") {
				matchedEntity = entity.Name
				isListCommand = true
				break
			}
		}
	} else {
		// Create command
		for _, entity := range availableEntities {
			if entity.Name == "run" && (strings.HasPrefix(lowerCmd, "run ") || lowerCmd == "run") {
				matchedEntity = "run"
				if len(command) > 3 {
					initialArg = strings.TrimSpace(command[3:])
				}
				break
			}

			if strings.HasPrefix(lowerCmd, "create "+entity.Name) ||
				strings.HasPrefix(lowerCmd, "create "+entity.Name+"s") ||
				strings.HasPrefix(lowerCmd, "make "+entity.Name) ||
				strings.HasPrefix(lowerCmd, "execute "+entity.Name) {
				matchedEntity = entity.Name
				// Extract arg
				prefixes := []string{"create ", "make ", "execute "}
				var prefix string
				for _, p := range prefixes {
					if strings.HasPrefix(lowerCmd, p+entity.Name) {
						prefix = p + entity.Name
						break
					}
					if strings.HasPrefix(lowerCmd, p+entity.Name+"s") {
						prefix = p + entity.Name + "s"
						break
					}
				}
				if len(command) > len(prefix) {
					initialArg = strings.TrimSpace(command[len(prefix):])
				}
				break
			}
		}
	}

	if matchedEntity != "" {
		if isListCommand {
			result, err := generateListQuery(matchedEntity, dbDriver, db)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		// Create
		return runConversationStep(command, state, dbDriver, userID, matchedEntity, initialArg, db)
	}

	// Triggers
	for _, entity := range availableEntities {
		flow, err := getQuestionFlow(db, entity.Name)
		if err == nil && len(flow.Triggers) > 0 {
			for _, trigger := range flow.Triggers {
				if strings.Contains(lowerCmd, strings.ToLower(trigger)) {
					// Found trigger, force new conversation
					return runConversationStep(command, AIConversationState{}, dbDriver, userID, entity.Name, "", db)
				}
			}
		}
	}
	
	// Unknown
	return map[string]interface{}{
		"execute": "false",
		"message": "I didn't understand that command. Try 'create menu', 'list menus', 'get data [table]' or 'help'.",
		"user_id": userID,
	}, nil
}


func runConversationStep(command string, state AIConversationState, dbDriver, userID, matchedEntity, initialArg string, db *gorm.DB) (map[string]interface{}, error) {
	lowerCmd := strings.ToLower(command)

	if lowerCmd == "exit" || lowerCmd == "cancel" {
		return map[string]interface{}{
			"execute":            "false",
			"message":            "Conversation cancelled.",
			"cancelled":          true,
			"conversation_state": "{}",
			"user_id":            userID,
		}, nil
	}

	// Set entity type
	if state.EntityType == "" {
		if matchedEntity != "" {
			state.EntityType = matchedEntity
		}
	}

	flow, err := getQuestionFlow(db, state.EntityType)
	if err != nil {
		return nil, err
	}

	// Handle Initial Arg for first Q
	if state.CurrentStep == 0 && initialArg != "" && len(flow.Questions) > 0 {
		firstQ := flow.Questions[0]
		state.CollectedData[firstQ.Key] = initialArg
		state.CurrentStep = 1
	}

	// Waiting Confirmation Phase
	if state.WaitingConfirmation {
		if lowerCmd == "execute" || lowerCmd == "yes" || lowerCmd == "confirm" {
			state.IsComplete = true
			query, params, reply, metaAction, err := generateQueryFromData(state.EntityType, state.CollectedData, dbDriver, userID, db)
			if err != nil {
				return nil, err
			}

			result := map[string]interface{}{
				"execute":            "true",
				"query":              query,
				"parameters":         params,
				"meta_action":        metaAction,
				"message":            reply,
				"completed":          true,
				"conversation_state": "{}",
				"user_id":            userID,
			}
			// Unmarshal query if JSON
			if strings.HasPrefix(strings.TrimSpace(query), "{") {
				var queryMap map[string]interface{}
				if err := json.Unmarshal([]byte(query), &queryMap); err == nil {
					for k, v := range queryMap {
						result[k] = v
					}
				}
			}
			return result, nil

		} else if lowerCmd == "review" {
			stateBytes, _ := json.Marshal(state)
			reviewText := buildDetailedReview(state.EntityType, state.CollectedData, flow)
			return map[string]interface{}{
				"execute":              "false",
				"message":              fmt.Sprintf("%s\n\nType 'execute' to proceed or 'cancel' to abort.", reviewText),
				"conversation_state":   string(stateBytes),
				"waiting_confirmation": true,
				"user_id":              userID,
			}, nil
		} else if lowerCmd == "cancel" || lowerCmd == "no" {
			return map[string]interface{}{
				"execute":   "false",
				"message":   "Operation cancelled.",
				"cancelled": true,
				"user_id":   userID,
			}, nil
		} else {
			stateBytes, _ := json.Marshal(state)
			summary := buildSummary(db, state.EntityType, state.CollectedData)
			return map[string]interface{}{
				"execute":              "false",
				"message":              fmt.Sprintf("Invalid response. Please type:\n• 'execute' to proceed\n• 'review' to see all questions and answers\n• 'cancel' to abort\n\n%s", summary),
				"conversation_state":   string(stateBytes),
				"waiting_confirmation": true,
				"user_id":              userID,
			}, nil
		}
	}

	// Capture answer from previous step if relevant
	// NOTE: If we just started (CurrentStep=0 -> 1 via initialArg), we don't capture `command` as answer.
	// If we are continuing (step > 0), `command` is the answer to step-1.
	// Logic matching main.go:
	if state.CurrentStep > 0 && state.CurrentStep <= len(flow.Questions) {
		// Only capture if we didn't just jump here via initialArg in this same call
		// But Wait: runConversationStep is called with `command`.
		// If matchedEntity != "" (New conversation), we handled initialArg.
		// If matchedEntity == "" (Continuing), `command` IS the answer.
		if matchedEntity == "" {
			prevQuestion := flow.Questions[state.CurrentStep-1]
			state.CollectedData[prevQuestion.Key] = command
		}
	}

	// Ask next question
	if state.CurrentStep < len(flow.Questions) {
		question := flow.Questions[state.CurrentStep]
		state.CurrentStep++

		stateBytes, _ := json.Marshal(state)
		response := map[string]interface{}{
			"execute":            "false",
			"message":            question.Text,
			"conversation_state": string(stateBytes),
			"question_key":       question.Key,
			"question_help":      question.Description,
			"progress":           fmt.Sprintf("Question %d of %d", state.CurrentStep, len(flow.Questions)),
			"user_id":            userID,
		}

		if question.OptionsQuery != "" && db != nil {
			var dynamicOptions []string
			if err := db.Raw(question.OptionsQuery).Scan(&dynamicOptions).Error; err == nil && len(dynamicOptions) > 0 {
				question.Options = dynamicOptions
			}
		}
		if len(question.Options) > 0 {
			response["options"] = question.Options
		}
		return response, nil
	}

	// All questions answered
	// Check single-step auto-execute
	/* if len(flow.Questions) == 1 && len(state.CollectedData) == 1 {
		state.IsComplete = true
		query, params, reply, metaAction, err := generateQueryFromData(state.EntityType, state.CollectedData, dbDriver, userID)
		if err != nil {
			return nil, err
		}
		result := map[string]interface{}{
			"execute":            "true",
			"query":              query,
			"parameters":         params,
			"meta_action":        metaAction,
			"message":            reply,
			"completed":          true,
			"conversation_state": "{}",
			"user_id":            userID,
		}
		if strings.HasPrefix(strings.TrimSpace(query), "{") {
			var queryMap map[string]interface{}
			if err := json.Unmarshal([]byte(query), &queryMap); err == nil {
				for k, v := range queryMap {
					result[k] = v
				}
			}
		}
		return result, nil
	} */

	// Summary and Confirm
	state.WaitingConfirmation = true
	stateBytes, _ := json.Marshal(state)
	summary := buildSummary(db, state.EntityType, state.CollectedData)

	return map[string]interface{}{
		"execute":              "false",
		"message":              fmt.Sprintf("✅ All information collected!\n\n%s\n\nType 'execute' to create this %s, or 'cancel' to abort.", summary, state.EntityType),
		"conversation_state":   string(stateBytes),
		"waiting_confirmation": true,
		"summary":              summary,
		"user_id":              userID,
	}, nil
}

// -- Helpers --

// AICommand represents the database structure for AI commands
type AICommand struct {
	AICommandID    int    `gorm:"primaryKey;column:aicommandid"`
	Name           string `gorm:"column:name"`
	Description    string `gorm:"column:description"`
	Questions      string `gorm:"column:questions"`
	Queries        string `gorm:"column:queries"`
	Params         string `gorm:"column:params"`
	Defaults       string `gorm:"column:defaults"`
	MetaAction     string `gorm:"column:metaaction"`
	SuccessMessage string `gorm:"column:successmessage"`
	Triggers       string `gorm:"column:triggers"`
}

// TableName overrides the table name used by User to `aicommand`
func (AICommand) TableName() string {
	return "aicommand"
}

func getQuestionFlow(db *gorm.DB, entityType string) (*AIQuestionFlow, error) {
	var cmd AICommand
	if err := db.Where("name = ?", entityType).First(&cmd).Error; err != nil {
		return nil, fmt.Errorf("command '%s' not found: %v", entityType, err)
	}

	var flow AIQuestionFlow
	flow.Description = cmd.Description
	flow.MetaAction = cmd.MetaAction
	flow.SuccessMessage = cmd.SuccessMessage

	// Unmarshal JSON fields
	json.Unmarshal([]byte(cmd.Questions), &flow.Questions)
	json.Unmarshal([]byte(cmd.Queries), &flow.Queries)
	// Handle ListQueries? Struct doesn't have it explicitly, maybe inside Queries or separate?
	// The DB schema provided didn't have listqueries column. 
	// Assuming logic needs adaptation or it's part of queries?
	// For now let's Initialize map
	if flow.Queries == nil { flow.Queries = make(map[string]string) }
	
	json.Unmarshal([]byte(cmd.Params), &flow.Params)
	json.Unmarshal([]byte(cmd.Defaults), &flow.Defaults)
	json.Unmarshal([]byte(cmd.Triggers), &flow.Triggers)

	// Defaults for safety
	if flow.Defaults == nil { flow.Defaults = make(map[string]string) }

	return &flow, nil
}

func getAvailableEntities(db *gorm.DB) []AIEntityInfo {
	var commands []AICommand
	if err := db.Find(&commands).Error; err != nil {
		fmt.Printf("[CompAI] Error listing commands: %v\n", err)
		return []AIEntityInfo{}
	}

	var entities []AIEntityInfo
	for _, cmd := range commands {
		description := fmt.Sprintf("Create a %s", cmd.Name)
		if cmd.Description != "" {
			description = cmd.Description
		}
		
		var triggers []string
		_ = json.Unmarshal([]byte(cmd.Triggers), &triggers)
		triggerStr := strings.Join(triggers, ", ")
		if triggerStr == "" { triggerStr = cmd.Name }

		entities = append(entities, AIEntityInfo{Name: cmd.Name, Description: description, Triggers: triggerStr})
	}
	return entities
}

func buildSummary(db *gorm.DB, entityType string, data map[string]string) string {
	var summary strings.Builder
	summary.WriteString("📋 Summary:\n")

	// Generic dumper for fields to avoid hardcoding switch for every type
	// But keeping switch for nice formatting if needed.
	// For now, let's iterate to be valid for ANY json file.
	// But order is random in map.
	// Let's use the flow to get order?
	flow, err := getQuestionFlow(db, entityType)
	if err == nil {
		for _, q := range flow.Questions {
			val := data[q.Key]
			if val != "" {
				summary.WriteString(fmt.Sprintf("  • %s: %s\n", q.Text, val))
			}
		}
	} else {
		for k, v := range data {
			summary.WriteString(fmt.Sprintf("  • %s: %s\n", k, v))
		}
	}
	return summary.String()
}

func buildDetailedReview(entityType string, data map[string]string, flow *AIQuestionFlow) string {
	var review strings.Builder
	review.WriteString("📝 Detailed Review - All Questions and Answers:\n\n")

	for i, question := range flow.Questions {
		answer := data[question.Key]
		if answer == "" {
			answer = "(not provided)"
		}
		review.WriteString(fmt.Sprintf("%d. %s\n", i+1, question.Text))
		review.WriteString(fmt.Sprintf("   Answer: %s\n\n", answer))
	}
	return review.String()
}

func generateQueryFromData(entityType string, data map[string]string, dbDriver, userID string, db *gorm.DB) (string, string, string, string, error) {
	flow, err := getQuestionFlow(db, entityType)
	if err != nil {
		return "", "", "", "", err
	}

	// Apply defaults
	for k, v := range flow.Defaults {
		if _, exists := data[k]; !exists || data[k] == "" {
			data[k] = processTemplate(v, data)
		}
	}

	// Helper for specific hardcoded things (restore if needed, or rely on flexible JSONs)
	if data["menuname"] != "" &&FuncMapHas(data, "menuname_slug") == false {
		data["menuname_slug"] = strings.ToLower(strings.ReplaceAll(data["menuname"], " ", "_"))
	}
	
	if entityType == "table" && data["columns"] != "" {
		processTableColumns(data, dbDriver)
	}

	queryTmpl, ok := flow.Queries[dbDriver]
	if !ok {
		if val, ok := flow.Queries["mysql"]; ok {
			queryTmpl = val
		} else if val, ok := flow.Queries["default"]; ok {
			queryTmpl = val
		} else {
			return "", "", "", "", fmt.Errorf("no query template for driver: %s", dbDriver)
		}
	}

	query := processTemplate(queryTmpl, data)

	var params []interface{}
	for _, paramName := range flow.Params {
		params = append(params, data[paramName])
	}
	jsonParams, _ := json.Marshal(params)

	successMsg := processTemplate(flow.SuccessMessage, data)
	if successMsg == "" {
		successMsg = fmt.Sprintf("%s created successfully!", strings.Title(entityType))
	}

	metaAction := flow.MetaAction
	if metaAction == "" {
		metaAction = "insert"
	}

	return query, string(jsonParams), successMsg, metaAction, nil
}

func processTemplate(tmpl string, data map[string]string) string {
	t, err := template.New("query").Parse(tmpl)
	if err != nil {
		return tmpl
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return tmpl
	}
	return buf.String()
}

func generateListQuery(entityType, dbDriver string, db *gorm.DB) (map[string]interface{}, error) {
	flow, err := getQuestionFlow(db, entityType)
	if err != nil {
		return nil, err
	}

	queryTmpl, ok := flow.ListQueries[dbDriver]
	if !ok {
		if val, ok := flow.ListQueries["mysql"]; ok {
			queryTmpl = val
		} else if val, ok := flow.ListQueries["default"]; ok {
			queryTmpl = val
		} else {
			return nil, fmt.Errorf("no list query defined for %s", entityType)
		}
	}

	query := processTemplate(queryTmpl, map[string]string{})

	return map[string]interface{}{
		"query":       query,
		"parameters":  "[]",
		"meta_action": "select",
	}, nil
}

func processTableColumns(data map[string]string, driver string) {
	columnsStr := data["columns"]
	primaryKey := data["primarykey"]
	columnDefs := []string{}
	columnParts := strings.Split(columnsStr, ",")

	for i, col := range columnParts {
		col = strings.TrimSpace(col)
		parts := strings.Split(col, ":")
		if len(parts) != 2 {
			continue
		}
		colName := strings.TrimSpace(parts[0])
		colType := strings.TrimSpace(parts[1])

		if primaryKey == "" && i == 0 {
			if strings.Contains(driver, "mysql") || strings.Contains(driver, "mariadb") {
				columnDefs = append(columnDefs, fmt.Sprintf("%s %s PRIMARY KEY AUTO_INCREMENT", colName, colType))
			} else if strings.Contains(driver, "postgres") {
				columnDefs = append(columnDefs, fmt.Sprintf("%s SERIAL PRIMARY KEY", colName))
			} else if strings.Contains(driver, "sqlite") {
				columnDefs = append(columnDefs, fmt.Sprintf("%s %s PRIMARY KEY AUTOINCREMENT", colName, colType))
			} else {
				columnDefs = append(columnDefs, fmt.Sprintf("%s %s PRIMARY KEY", colName, colType))
			}
		} else {
			columnDefs = append(columnDefs, fmt.Sprintf("%s %s", colName, colType))
		}
	}
	data["column_defs"] = strings.Join(columnDefs, ", ")
}

func handleGetData(db *gorm.DB, tableName string) (interface{}, error) {
	var results []map[string]interface{}
	// Safety check: tableName should only contain alphanumeric/underscore to prevent injection
	// Or use GORM safely? GORM raw with table name param doesn't always work for FROM clause
	// But for simple "get data", maybe we just trust input if internally authenticated?
	// Let's do basic sanitization
	cleanTable := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, tableName)
	
	if cleanTable == "" { return nil, fmt.Errorf("invalid table name") }

	// Limit to 50 for safety
	if err := db.Table(cleanTable).Limit(50).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func FuncMapHas(m map[string]string, key string) bool {
    _, ok := m[key]
    return ok
}

func formatDataAsTable(data []map[string]interface{}) string {
	if len(data) == 0 {
		return "No data found."
	}
	// Collect all keys
	keys := make(map[string]bool)
	var header []string
	for _, row := range data {
		for k := range row {
			if !keys[k] {
				keys[k] = true
				header = append(header, k)
			}
		}
	}
	sort.Strings(header)

	var b strings.Builder
	// Header
	for _, h := range header {
		b.WriteString("| " + h + " ")
	}
	b.WriteString("|\n")
	// Separator
	for range header {
		b.WriteString("| --- ")
	}
	b.WriteString("|\n")
	// Rows
	for _, row := range data {
		for _, h := range header {
			val := ""
			if v, ok := row[h]; ok {
				val = fmt.Sprintf("%v", v)
			}
			// Sanitize newlines
			val = strings.ReplaceAll(val, "\n", " ")
			b.WriteString("| " + val + " ")
		}
		b.WriteString("|\n")
	}
	return b.String()
}

// generateExcelFromData creates an Excel file from query results
func generateExcelFromData(tableName string, data []map[string]interface{}) (string, error) {
	// Import excelize at the top of file if not already imported
	excel := excelize.NewFile()
	defer excel.Close()
	
	sheetName := "Sheet1"
	excel.SetSheetName(sheetName, tableName)
	
	if len(data) == 0 {
		return "", fmt.Errorf("no data to export")
	}
	
	// Get headers from first row
	var headers []string
	headerMap := make(map[string]bool)
	for _, row := range data {
		for key := range row {
			if !headerMap[key] {
				headers = append(headers, key)
				headerMap[key] = true
			}
		}
	}
	sort.Strings(headers)
	
	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		excel.SetCellValue(tableName, cell, header)
	}
	
	// Write data rows
	for rowIdx, row := range data {
		for colIdx, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			value := row[header]
			excel.SetCellValue(tableName, cell, value)
		}
	}
	
	// Auto-fit columns
	for i := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		excel.SetColWidth(tableName, col, col, 15)
	}
	
	// Save to temp file
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.xlsx", tableName, timestamp)
	filePath := filepath.Join(os.TempDir(), fileName)
	
	if err := excel.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("failed to save Excel file: %v", err)
	}
	
	return filePath, nil
}
