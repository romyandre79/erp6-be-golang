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

	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// -- Structs converted from component-ai/main.go --

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
		// Config params
		token    string
		provider string
		model    string
		baseURL  string
		// Document context params
		documentIDs      string
		useAllDocuments  bool
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))

		// Fallback for legacy params (command, user_id)
		// Only check InputName if val is empty AND it's a legacy param that might be passed via context variable with same name
		// For explicit config params (provider, token), we should NOT fallback to the key name.

		switch strings.ToLower(p.InputName) {
		case "command", "message", "text", "input":
			if val == "" {
				val = ResolveParam(ctx.FiberCtx, p.InputName)
			}
			if !strings.HasPrefix(val, "$") {
				command = val
			}

		case "user_id":
			if val == "" {
				val = ResolveParam(ctx.FiberCtx, p.InputName)
			}
			if !strings.HasPrefix(val, "$") {
				userID = val
			}

		case "conversation_state":
			if val == "" {
				val = ResolveParam(ctx.FiberCtx, p.InputName)
			}
			if !strings.HasPrefix(val, "$") {
				conversationState = val
			}

		case "token", "api_key", "apikey":
			if val != "" {
				token = val
			}

		case "provider", "llm_provider":
			if val != "" {
				provider = val
			}

		case "model", "llm_model":
			if val != "" {
				model = val
			}

		case "base_url", "baseurl":
			if val != "" {
				baseURL = val
			}

		case "document_ids", "documentids":
			if val != "" {
				documentIDs = val
			}

		case "use_all_documents", "usealldocuments":
			if strings.ToLower(val) == "true" {
				useAllDocuments = true
			}
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

	aiConfig := models.AIConfig{
		Token:    token,
		Provider: provider,
		Model:    model,
		BaseURL:  baseURL,
	}

	fmt.Printf("[CompAI DEBUG] Config Parsed -> Provider: '%s', Model: '%s', Token Set: %v\n", provider, model, token != "")

	// Pass document context to processAI
	result, err := processAI(command, conversationState, dbDriver, userID, ctx.DB, aiConfig, documentIDs, useAllDocuments)
	if err != nil {
		return err
	}

	// Inject results into Extras for subsequent components
	if ctx.Extras != nil {
		for k, v := range result {
			ctx.Extras[k] = v
		}
	} else {
		// Initialize if nil (though usually initialized by engine)
		ctx.Extras = make(map[string]interface{})
		for k, v := range result {
			ctx.Extras[k] = v
		}
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

	fmt.Printf("[ExecuteAIResult T_RACE] Started. Result: %+v\n", aiResult)

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

			// Save PARENT state to prevent pollution by Child
			parentWfEngine := ctx.FiberCtx.Locals("wfEngine")
			parentComponents := ctx.FiberCtx.Locals("components")
			parentTerminated := ctx.FiberCtx.Locals("flowTerminated")
			parentNested := ctx.FiberCtx.Locals("nestedWorkflow")

			// Set Nested Mode for Child
			ctx.FiberCtx.Locals("nestedWorkflow", true)

			fmt.Printf("[ExecuteAIResult] Running sub-workflow '%s' (Nested=true)\n", workflowName)
			err := ExecuteFlow(ctx.FiberCtx, db, workflowName, false, params)

			// Restore PARENT state
			ctx.FiberCtx.Locals("wfEngine", parentWfEngine)
			ctx.FiberCtx.Locals("components", parentComponents)
			ctx.FiberCtx.Locals("flowTerminated", parentTerminated)
			ctx.FiberCtx.Locals("nestedWorkflow", parentNested)

			return err
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
		// Map known actions to Web Scraper
		switch action {
		case "extract_one_data", "extract_data", "extract_links", "extract_text", "extract_images",
			"submit_and_extract", "extract_with_regex", "solve_captcha", "extract_hierarchy", "get_html":
			componentName = "Web Scraper"
		}
	}

	// Check if this is a workflow-routing command (has meta_action but no action)
	if componentName == "" {
		if metaAction, ok := aiResult["meta_action"].(string); ok && metaAction != "" {
			// This is a workflow-routing command (e.g., "Data Customer")
			// For web: workflow executed by Decision component
			// For WhatsApp: execute the component directly here
			fmt.Printf("[ExecuteAIResult] Workflow-routing command detected: meta_action=%s\n", metaAction)

			// Check if this is WhatsApp context (no FiberCtx)
			if ctx.FiberCtx == nil {
				fmt.Printf("[ExecuteAIResult] WhatsApp context - executing component directly\n")

				// For WhatsApp, execute the Search component directly or handle specific meta_actions
				// Map meta_action to component and parameters
				switch metaAction {
				case "data_customer":
					// Execute Search component with customer table
					// Query directly and format
					var customerData []map[string]interface{}
					if err := db.Table("customer").Find(&customerData).Error; err != nil {
						return fmt.Errorf("failed to get customer data: %v", err)
					}

					// Format as table and store in aiResult for WhatsApp handler to use
					formattedMsg := helpers.FormatDataAsTable(customerData)
					aiResult["message"] = formattedMsg
					fmt.Printf("[ExecuteAIResult] Formatted message for WhatsApp: %s\n", formattedMsg)

				default:
					// Check if we can map other meta_actions
					// If URL is present, fallback to Web Scraper
					if url, ok := aiResult["url"].(string); ok && url != "" {
						fmt.Printf("[ExecuteAIResult] URL found, falling back to Web Scraper for meta_action: %s\n", metaAction)
						componentName = "Web Scraper"
						// Force browser method for consistency
						aiResult["method"] = "browser"
					} else {
						// If not implemented, just return nil (Workflow routing will happen in Decision node for Web)
						// For WA, this means no action, but maybe message is already set
						fmt.Printf("[ExecuteAIResult] No direct WA handler for meta_action: %s\n", metaAction)
					}
				}
			}

			// If we set a componentName (e.g. fallback), don't return nil, let it proceed to execution
			if componentName == "" {
				return nil
			}
		}
	}

	if componentName == "" {
		return fmt.Errorf("no component name in AI result")
	}

	// Execute component
	handler, ok := GetComponent(componentName)
	if !ok {
		fmt.Printf("[ExecuteAIResult TRACE] Component '%s' not found!\n", componentName)
		return fmt.Errorf("component '%s' not found", componentName)
	}

	fmt.Printf("[ExecuteAIResult TRACE] Executing component: %s\n", componentName)
	err := handler.Execute(ctx)
	fmt.Printf("[ExecuteAIResult TRACE] Component execution finished. Error: %v\n", err)
	return err
}

func processAI(command, stateJSON, dbDriver, userID string, db *gorm.DB, config models.AIConfig, documentIDs string, useAllDocuments bool) (map[string]interface{}, error) {
	var state models.AIConversationState

	// Parse existing state or create new one
	// Parse existing state or create new one
	if stateJSON != "" && stateJSON != "{}" {
		if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
			state = models.AIConversationState{CollectedData: make(map[string]string)}
		} else {
			// Ensure map is initialized if json had null
			if state.CollectedData == nil {
				state.CollectedData = make(map[string]string)
			}
		}
	} else {
		state = models.AIConversationState{CollectedData: make(map[string]string), History: []models.AIMessage{}}
	}

	// Add user message to history
	state.History = append(state.History, models.AIMessage{
		Role:      "user",
		Content:   command,
		Timestamp: time.Now().Unix(),
	})

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

	// Continue existing conversation
	if state.EntityType != "" {
		return runConversationStep(command, state, dbDriver, userID, "", "", db)
	}

	// New Conversation Logic
	availableEntities := getAvailableEntities(db)
	matchedEntity := ""
	initialArg := ""
	isListCommand := false

	/*if strings.HasPrefix(lowerCmd, "list ") {
		for _, entity := range availableEntities {
			if strings.HasPrefix(lowerCmd, "list "+entity.Name) || strings.HasPrefix(lowerCmd, "list "+entity.Name+"s") {
				matchedEntity = entity.Name
				isListCommand = true
				break
			}
		}
	} else {*/
	// Create command
	for _, entity := range availableEntities {
		if entity.Name == "run" && (strings.HasPrefix(lowerCmd, "run ") || lowerCmd == "run") {
			matchedEntity = "run"
			if len(command) > 3 {
				initialArg = strings.TrimSpace(command[3:])
			}
			break
		}

		if strings.HasPrefix(lowerCmd, "execute "+entity.Name) {
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
	//}

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
					return runConversationStep(command, models.AIConversationState{}, dbDriver, userID, entity.Name, "", db)
				}
			}
		}
	}

	// Unknown command - Delegate to LLM
	fmt.Printf("[CompAI] No strict command matched. Delegating to LLM...\n")
	return delegateToLLM(command, state, dbDriver, userID, db, config, documentIDs, useAllDocuments)
}

func delegateToLLM(command string, state models.AIConversationState, driver, userID string, db *gorm.DB, config models.AIConfig, documentIDs string, useAllDocuments bool) (map[string]interface{}, error) {
	// 1. Get Schema Context
	schemaTables, err := ReverseEngineerDatabase(db)
	if err != nil {
		fmt.Printf("[CompAI] Warning: Failed to get schema: %v\n", err)
	}

	// Simplify schema for prompt
	var schemaSummary strings.Builder
	schemaSummary.WriteString("Database Schema:\n")
	for _, t := range schemaTables {
		schemaSummary.WriteString(fmt.Sprintf("- Table: %s\n", t.Name))
		for _, c := range t.Columns {
			schemaSummary.WriteString(fmt.Sprintf("  - %s (%s)\n", c.Name, c.Type))
		}
	}

	// 1.5. Get Document Context
	var documentContext strings.Builder
	if documentIDs != "" || useAllDocuments {
		var documents []models.Document
		var err error

		if useAllDocuments {
			// Get all documents for this user
			err = db.Where("userid = ?", userID).Find(&documents).Error
			fmt.Printf("[CompAI] Loading all documents for user %s\n", userID)
		} else {
			// Get specific documents by IDs
			ids := strings.Split(documentIDs, ",")
			var cleanIDs []string
			for _, id := range ids {
				if trimmed := strings.TrimSpace(id); trimmed != "" {
					cleanIDs = append(cleanIDs, trimmed)
				}
			}
			err = db.Where("documentid IN (?) AND userid = ?", cleanIDs, userID).Find(&documents).Error
			fmt.Printf("[CompAI] Loading documents with IDs: %v for user %s\n", cleanIDs, userID)
		}

		if err != nil {
			fmt.Printf("[CompAI] Warning: Failed to load documents: %v\n", err)
		} else if len(documents) > 0 {
			documentContext.WriteString("\n\nAvailable Documents:\n")
			for _, doc := range documents {
				documentContext.WriteString(fmt.Sprintf("\n[Document ID: %d - %s]\n", doc.DocumentID, doc.FileName))
				documentContext.WriteString(fmt.Sprintf("Type: %s | Size: %d bytes\n", doc.FileType, doc.FileSize))
				documentContext.WriteString("Content:\n")
				// Limit document content to avoid token overflow (e.g., 5000 chars per doc)
				maxChars := 5000
				if len(doc.ExtractedText) > maxChars {
					documentContext.WriteString(doc.ExtractedText[:maxChars])
					documentContext.WriteString("\n... (truncated)\n")
				} else {
					documentContext.WriteString(doc.ExtractedText)
				}
				documentContext.WriteString("\n---\n")
			}
			fmt.Printf("[CompAI] Loaded %d documents into context\n", len(documents))
		}
	}

	// 2. Build Prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are a helpful database assistant for an ERP system. ")
	promptBuilder.WriteString("You have access to the following database schema:\n")
	promptBuilder.WriteString(schemaSummary.String())
	
	// Add document context if available
	if documentContext.Len() > 0 {
		promptBuilder.WriteString(documentContext.String())
		promptBuilder.WriteString("\nYou can reference information from these documents when answering questions.\n")
	}
	
	promptBuilder.WriteString("\n\nAnswer the user's question. If you need to query the database, output the SQL query in valid JSON format like: {\"action\": \"query\", \"sql\": \"SELECT ...\"}. \n")
	promptBuilder.WriteString("If you can answer without querying (or have the result), just provide the answer.\n\n")

	// Add History
	for _, msg := range state.History {
		role := "User"
		if msg.Role == "assistant" {
			role = "Assistant"
		}
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}
	promptBuilder.WriteString("Assistant:")

	fullPrompt := promptBuilder.String()

	// 3. Call LLM (Using ENV for config for now, or default)

	// 3. Call LLM (Loop for Agentic behavior - max 5 turns)
	maxTurns := 5

	// Prioritize Params -> Env
	// Determine provider first, as it influences token lookup
	provider := config.Provider
	if provider == "" {
		provider = os.Getenv("LLM_PROVIDER")
	}

	token := config.Token
	if token == "" {
		if strings.ToLower(provider) == "gemini" {
			token = os.Getenv("GEMINI_API_KEY")
		} else if strings.EqualFold(provider, "claude") || strings.EqualFold(provider, "anthropic") {
			token = os.Getenv("ANTHROPIC_API_KEY")
		}

		if token == "" {
			token = os.Getenv("OPENAI_API_KEY")
		}
	}

	model := config.Model
	if model == "" {
		model = os.Getenv("LLM_MODEL")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		// Check for Ollama-specific env var first
		baseURL = os.Getenv("OLLAMA_BASE_URL")
		// Fallback to generic LLM_BASE_URL
		if baseURL == "" {
			baseURL = os.Getenv("LLM_BASE_URL")
		}
	}

	// Fallback if no env
	if token == "" {
		return map[string]interface{}{
			"message": "I'm sorry, I can't help with that yet (LLM Token not configured).",
			"user_id": userID,
		}, nil
	}

	fmt.Printf("[CompAI DEBUG] delegateToLLM -> Provider: '%s', Model: '%s'\n", provider, model)

	var finalResponse string

	for i := 0; i < maxTurns; i++ {
		// Re-build prompt with latest history
		var currentPromptBuilder strings.Builder
		currentPromptBuilder.WriteString(fullPrompt) // Base prompt

		// Add dynamic history (including tool outputs from this session)
		// Note: state.History is the *conversation* history.
		// We need to handle the *internal* reasoning loop history.
		// For simplicity, let's just append internal turns to state.History temporarily?
		// Better: Maintain a local conversation buffer for this turn, or just rely on state.History if we commit to it.
		// Let's commit to state.History as it allows debugging.

		// Wait, 'fullPrompt' already has the history up to start of function.
		// We actually need to re-generate the history part of the prompt in each loop iteration
		// OR just append the new messages to the prompt.
		// Valid approach: Just re-generate prompt from state.History

		var loopPromptBuilder strings.Builder
		loopPromptBuilder.WriteString("You are a helpful database assistant for an ERP system. ")
		loopPromptBuilder.WriteString("You have access to the following database schema:\n")
		loopPromptBuilder.WriteString(schemaSummary.String())
		loopPromptBuilder.WriteString("\n\nAnswer the user's question. If you need to query the database, output the SQL query in valid JSON format like: {\"action\": \"query\", \"sql\": \"SELECT ...\"}. \n")
		loopPromptBuilder.WriteString("If you can answer without querying (or have the result), just provide the answer.\n\n")

		for _, msg := range state.History {
			role := "User"
			if msg.Role == "assistant" {
				role = "Assistant"
			}
			if msg.Role == "system" {
				role = "System"
			} // For Tool Outputs
			loopPromptBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
		}
		loopPromptBuilder.WriteString("Assistant:")
		currentPrompt := loopPromptBuilder.String()

		// Prepare candidate providers for rotation/fallback
		var candidates []models.AIConfig

		// Parse Comma-Separated Configurations - we do this manually below via split

		// If Env var was used and params were empty, respect that (already handled because we pass config.* which might be from Params)
		// Wait, earlier logic set 'provider' var from logic: params -> env.
		// But here we want to re-evaluate based on the potentially raw inputs or just use the vars passed in?
		// The `config` struct passed to delegateToLLM contains the raw values from ctx.Params (or defaults if we logic'd them?)
		// Actually processAI receives `config AIConfig` which is constructed in handleAI directly from params.
		// So `config.Provider` is exactly what came from params.
		// BUT, if params were empty, we want to start with the ENV var value.

		// Let's re-resolve the "Main" strings to iterate on.
		pStr := config.Provider
		if pStr == "" {
			pStr = os.Getenv("LLM_PROVIDER")
		}

		tStr := config.Token
		// If token param is empty, we don't automatically grab one single ENNV because it depends on the provider list.
		// But if they provided a single provider in env and no token param, we might want to grab the matching env key.
		// We'll handle "missing token" inside the loop by looking up ENV.

		mStr := config.Model
		if mStr == "" {
			mStr = os.Getenv("LLM_MODEL")
		}

		pParts := strings.Split(pStr, ",")
		tParts := strings.Split(tStr, ",")
		mParts := strings.Split(mStr, ",")

		trim := func(s []string) []string {
			var r []string
			for _, v := range s {
				if t := strings.TrimSpace(v); t != "" {
					r = append(r, t)
				}
			}
			return r
		}
		pParts = trim(pParts)
		// Don't trim empty tokens entirely? Well, "key1,,key3" -> middle one empty?
		// strings.Split gives empty strings. The trim function above removes them.
		// If user did "gemini,openai" and "key1," (missing second), we want to detect that.
		// Let's just create a safe accessor.

		// Helper to safely get item from slice or empty
		getAt := func(slice []string, i int) string {
			if i < len(slice) {
				return strings.TrimSpace(slice[i])
			}
			return ""
		}

		seenCandidates := make(map[string]bool) // Key: Provider+Model+Token

		// Helper to get models for a provider, placing the preferred ones first
		getModels := func(prov, preferredModelsStr string) []string {
			p := strings.ToLower(prov)
			var defaults []string
			if p == "gemini" {
				//gemini-3-flash-preview:gemini-2.5-pro:gemini-2.5-flash:gemini-2.5-flash-preview-09-2025:gemini-2.5-flash-lite:gemini-2.5-flash-lite-preview-09-2025:gemini-2.0-flash:gemini-2.0-flash-lite:gemini-1.5-flash:gemini-1.5-pro:gemini-1.0-pro
				defaults = []string{"gemini-3-flash-preview","gemini-2.5-pro","gemini-2.5-flash","gemini-2.5-flash-preview-09-2025","gemini-2.5-flash-lite","gemini-2.5-flash-lite-preview-09-2025","gemini-2.0-flash","gemini-2.0-flash-lite","gemini-1.5-flash", "gemini-1.5-pro", "gemini-1.0-pro"}
			} else if p == "openai" {
				defaults = []string{"gpt-4o", "gpt-4-turbo", "gpt-3.5-turbo"}
			} else if p == "claude" || p == "anthropic" {
				defaults = []string{"claude-3-5-sonnet-20240620", "claude-3-opus-20240229", "claude-3-haiku-20240307"}
			}

			var final []string
			seen := make(map[string]bool)

			// Parse user provided list
			if preferredModelsStr != "" {
				parts := strings.Split(preferredModelsStr, ",")
				for _, part := range parts {
					m := strings.TrimSpace(part)
					if m != "" && !seen[m] {
						final = append(final, m)
						seen[m] = true
					}
				}
			}

			// Append defaults if not seen
			for _, m := range defaults {
				if !seen[m] {
					final = append(final, m)
					seen[m] = true
				}
			}
			return final
		}

		// 1. Build candidates from explicit lists
		for i, rawProv := range pParts {
			provName := strings.ToLower(rawProv)

			// Get corresponding token
			userToken := getAt(tParts, i)

			// Resolve Token if empty mechanism
			if userToken == "" {
				if provName == "gemini" {
					userToken = os.Getenv("GEMINI_API_KEY")
				}
				if provName == "openai" {
					userToken = os.Getenv("OPENAI_API_KEY")
				}
				if provName == "claude" || provName == "anthropic" {
					userToken = os.Getenv("ANTHROPIC_API_KEY")
				}
			}

			if userToken == "" {
				fmt.Printf("[CompAI DEBUG] Skipping provider '%s' (index %d): No token found\n", provName, i)
				continue
			}

			// Resolve Models
			// Logic: Comma separates models within a provider chunk
			// Note: Colon is NOT a separator - it's part of model names (e.g., Ollama's "gemma3:1b")
			// Example: Provider="gemini,openai" Model="gemini-1.5-flash,gemini-1.5-pro, gpt-4o"

			targetModelStr := ""

			if len(pParts) == 1 {
				// Single Provider: Treat entire model string as the list for this provider
				targetModelStr = mStr
			} else {
				// Multi Provider: Get the chunk corresponding to this provider index
				modelChunk := getAt(mParts, i)

				// Handle mismatched lengths (e.g. 2 providers, 1 model string)
				// If missing, try to use the last available chunk or default
				if modelChunk == "" && len(mParts) > 0 {
					modelChunk = getAt(mParts, len(mParts)-1)
				}

				targetModelStr = modelChunk
			}

			modelList := getModels(provName, targetModelStr)
			for _, m := range modelList {
				key := provName + "|" + m + "|" + userToken
				if !seenCandidates[key] {
					candidates = append(candidates, models.AIConfig{
						Token:    userToken,
						Provider: provName,
						Model:    m,
						BaseURL:  baseURL,
					})
					seenCandidates[key] = true
				}
			}
		}

		// 2. Auto-Discover Fallbacks (Env) if not already explicitly added
		// (This covers the case where user didn't even put them in the list)

		addFallback := func(pName, envKey string) {
			token := os.Getenv(envKey)
			if token == "" {
				return
			}

			// Check if we already have this provider coverage?
			// Maybe checking "gemini" presence in pParts is enough?
			// But user might have "gemini" in pParts but with a specific token. We might want to add Env-based Gemini as backup?
			// Let's just add it. Duplication check via `seenCandidates` handles duplicates (same token/model).
			// If token is different (Env vs Param), it's a valid new candidate!

			modelList := getModels(pName, "")
			// Limit fallbacks to 2 to not spam
			if len(modelList) > 2 {
				modelList = modelList[:2]
			}

			for _, m := range modelList {
				key := pName + "|" + m + "|" + token
				if !seenCandidates[key] {
					candidates = append(candidates, models.AIConfig{
						Token:    token,
						Provider: pName,
						Model:    m,
						BaseURL:  baseURL,
					})
					seenCandidates[key] = true
				}
			}
		}

		addFallback("gemini", "GEMINI_API_KEY")
		addFallback("openai", "OPENAI_API_KEY")
		addFallback("claude", "ANTHROPIC_API_KEY")

		var err error
		var llmResponse string
		var lastError error
		success := false

		// Try each candidate
		for idx, cand := range candidates {
			pName := strings.TrimSpace(strings.ToLower(cand.Provider))
			fmt.Printf("[CompAI DEBUG] Attempt %d using %s (Model: %s)\n", idx+1, pName, cand.Model)

				if pName == "gemini" {
				llmResponse, err = RunGemini(cand.Token, cand.Model, currentPrompt, cand.BaseURL)
			} else if pName == "claude" || pName == "anthropic" {
				llmResponse, err = RunAnthropic(cand.Token, cand.Model, currentPrompt, cand.BaseURL)
			} else {
				// Default to OpenAI (also works for Ollama with custom baseURL)
				llmResponse, err = RunOpenAI(cand.Token, cand.BaseURL, cand.Model, currentPrompt)
			}

			if err == nil {
				success = true
				break // Success!
			} else {
				fmt.Printf("[CompAI DEBUG] Attempt %d failed: %v\n", idx+1, err)
				lastError = err
			}
		}

		if !success {
			return nil, fmt.Errorf("All LLM providers failed. Last error: %v", lastError)
		}

		fmt.Printf("[CompAI DEBUG] Turn %d LLM Response: %s\n", i, llmResponse)

		// Add Assistant response to history
		state.History = append(state.History, models.AIMessage{
			Role:      "assistant",
			Content:   llmResponse,
			Timestamp: time.Now().Unix(),
		})

		finalResponse = llmResponse

		// Check for JSON action
		// Sanitize markdown code blocks if present
		cleanResponse := strings.TrimSpace(llmResponse)
		cleanResponse = strings.ReplaceAll(cleanResponse, "```json", "")
		cleanResponse = strings.ReplaceAll(cleanResponse, "```", "")
		cleanResponse = strings.TrimSpace(cleanResponse)

		// Find JSON start/end just in case there is text around it
		jsonStart := strings.Index(cleanResponse, "{")
		jsonEnd := strings.LastIndex(cleanResponse, "}")

		if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
			possibleJSON := cleanResponse[jsonStart : jsonEnd+1]

			var action map[string]interface{}
			if err := json.Unmarshal([]byte(possibleJSON), &action); err == nil {
				if act, ok := action["action"].(string); ok && act == "query" {
					sqlQuery, _ := action["sql"].(string)
					fmt.Printf("[CompAI DEBUG] Executing SQL: %s\n", sqlQuery)

					// EXECUTE SQL
					var queryResult []map[string]interface{}
					if err := db.Raw(sqlQuery).Scan(&queryResult).Error; err != nil {
						// SQL Error
						toolOutput := fmt.Sprintf("SQL Error: %v", err)
						state.History = append(state.History, models.AIMessage{
							Role:      "system",
							Content:   toolOutput,
							Timestamp: time.Now().Unix(),
						})
						fmt.Printf("[CompAI DEBUG] Tool Output: %s\n", toolOutput)
					} else {
						resultBytes, _ := json.Marshal(queryResult)
						toolOutput := fmt.Sprintf("Query Result: %s", string(resultBytes))
						state.History = append(state.History, models.AIMessage{
							Role:      "system",
							Content:   toolOutput,
							Timestamp: time.Now().Unix(),
						})
						fmt.Printf("[CompAI DEBUG] Tool Output (Len): %d bytes\n", len(toolOutput))
					}
					// Continue loop to let LLM analyze result
					continue
				}
			}
		}

		// If no action found, or not a query, we are done
		break
	}

	stateBytes, _ := json.Marshal(state)
	
	// Try to parse finalResponse as JSON and extract "message" field if present
	cleanMessage := finalResponse
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(finalResponse), &jsonResponse); err == nil {
		// Successfully parsed as JSON
		if msg, ok := jsonResponse["message"].(string); ok && msg != "" {
			// Extract just the message field for cleaner output
			cleanMessage = msg
		}
	}
	
	return map[string]interface{}{
		"message":            cleanMessage,
		"conversation_state": string(stateBytes),
		"user_id":            userID,
	}, nil
}

func runConversationStep(command string, state models.AIConversationState, dbDriver, userID, matchedEntity, initialArg string, db *gorm.DB) (map[string]interface{}, error) {
	lowerCmd := strings.ToLower(command)

	if lowerCmd == "exit" || lowerCmd == "cancel" || lowerCmd == "batal" {
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
		switch lowerCmd {
		case "execute", "yes", "confirm", "lanjut":
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

			// Add collected data to result so it's available for parameters
			for k, v := range state.CollectedData {
				result[k] = v
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

		case "review", "kesimpulan", "summary":
			stateBytes, _ := json.Marshal(state)
			reviewText := buildDetailedReview(state.EntityType, state.CollectedData, flow)
			return map[string]interface{}{
				"execute":              "false",
				"message":              fmt.Sprintf("%s\n\nType 'execute' or 'lanjut' or 'confirm' or 'yes' to proceed or 'cancel' or 'no' or 'batal' to abort.", reviewText),
				"conversation_state":   string(stateBytes),
				"waiting_confirmation": true,
				"user_id":              userID,
			}, nil
		case "cancel", "no", "batal":
			return map[string]interface{}{
				"execute":   "false",
				"message":   "Operation cancelled.",
				"cancelled": true,
				"user_id":   userID,
			}, nil
		default:
			stateBytes, _ := json.Marshal(state)
			summary := buildSummary(db, state.EntityType, state.CollectedData)
			return map[string]interface{}{
				"execute":              "false",
				"message":              fmt.Sprintf("Invalid response. Please type:\n• 'execute' or 'lanjut' or 'confirm' or 'yes' to proceed\n• 'review' or 'kesimpulan' or 'summary' to see all questions and answers\n• 'cancel' or 'no' or 'batal' to abort\n\n%s", summary),
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

func getQuestionFlow(db *gorm.DB, entityType string) (*models.AIQuestionFlow, error) {
	var cmd models.AICommand
	if err := db.Where("name = ?", entityType).First(&cmd).Error; err != nil {
		return nil, fmt.Errorf("command '%s' not found: %v", entityType, err)
	}

	var flow models.AIQuestionFlow
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
	if flow.Queries == nil {
		flow.Queries = make(map[string]string)
	}

	json.Unmarshal([]byte(cmd.Params), &flow.Params)
	json.Unmarshal([]byte(cmd.Defaults), &flow.Defaults)
	json.Unmarshal([]byte(cmd.Triggers), &flow.Triggers)

	// Defaults for safety
	if flow.Defaults == nil {
		flow.Defaults = make(map[string]string)
	}

	return &flow, nil
}

func getAvailableEntities(db *gorm.DB) []models.AIEntityInfo {
	var commands []models.AICommand
	if err := db.Find(&commands).Error; err != nil {
		fmt.Printf("[CompAI] Error listing commands: %v\n", err)
		return []models.AIEntityInfo{}
	}

	var entities []models.AIEntityInfo
	for _, cmd := range commands {
		description := fmt.Sprintf("Create a %s", cmd.Name)
		if cmd.Description != "" {
			description = cmd.Description
		}

		var triggers []string
		_ = json.Unmarshal([]byte(cmd.Triggers), &triggers)
		triggerStr := strings.Join(triggers, ", ")
		if triggerStr == "" {
			triggerStr = cmd.Name
		}

		entities = append(entities, models.AIEntityInfo{Name: cmd.Name, Description: description, Triggers: triggerStr})
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

func buildDetailedReview(entityType string, data map[string]string, flow *models.AIQuestionFlow) string {
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

	// Format success message
	successMsg := processTemplate(flow.SuccessMessage, data)
	if successMsg == "" {
		successMsg = fmt.Sprintf("%s processed successfully!", strings.Title(entityType))
	}

	metaAction := flow.MetaAction
	if metaAction == "" {
		metaAction = "workflow"
	}

	// Determine query template based on driver
	queryTmpl := ""
	if val, ok := flow.Queries[dbDriver]; ok {
		queryTmpl = val
	} else if val, ok := flow.Queries["default"]; ok {
		queryTmpl = val
	}

	// Process query template (this contains the Scraper JSON config for scraping commands)
	query := processTemplate(queryTmpl, data)

	return query, "[]", successMsg, metaAction, nil
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
