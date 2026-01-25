package admin

import (
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	dbgenerator "erp6-be-golang/core/generator"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AiCommandRequest struct {
	Message string `json:"message"`
}

func AiCommandHandler(c *fiber.Ctx, db *gorm.DB) error {
	var body AiCommandRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	message := strings.TrimSpace(body.Message)
	lowerMsg := strings.ToLower(message)

	if strings.HasPrefix(lowerMsg, "create widget") {
		return handleCreateWidget(c, db, message)
	} else if strings.HasPrefix(lowerMsg, "create menu") {
		return handleCreateMenu(c, db, message)
	} else if strings.HasPrefix(lowerMsg, "create workflow") {
		return handleCreateWorkflow(c, db, message)
	}

	return helpers.SuccessResponse(c, "UNKNOWN_COMMAND", fiber.Map{
		"reply": "I'm sorry, I didn't understand that command. Try 'create widget [name]', 'create menu [name]', or 'create workflow [name]'.",
	})
}

func handleCreateWidget(c *fiber.Ctx, db *gorm.DB, message string) error {
	// message: "create widget WidgetName"
	parts := strings.Split(message, " ")
	if len(parts) < 3 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_COMMAND", "Please specify a widget name.")
	}
	widgetName := strings.Join(parts[2:], " ") // in case name has spaces? strict to one word? Let's allow spaces for title but usually name is unique.

	// Check if exists
	var count int64
	db.Model(&models.Widget{}).Where("widgetname = ?", widgetName).Count(&count)
	if count > 0 {
		return helpers.SuccessResponse(c, "WIDGET_EXISTS", fiber.Map{
			"reply": fmt.Sprintf("Widget '%s' already exists.", widgetName),
		})
	}

	newWidget := models.Widget{
		Widgetname:    widgetName,
		Widgettitle:   widgetName,
		Widgetversion: "1.0.0",
		Widgetby:      "AI Assistant",
		Description:   "Created via AI Assistant",
		Recordstatus:  1,
		Moduleid:      1, // Default module? Or ask? For now default.
	}

	if err := db.Create(&newWidget).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error())
	}

	return helpers.SuccessResponse(c, "WIDGET_CREATED", fiber.Map{
		"reply": fmt.Sprintf("Success! Widget '%s' has been created.", widgetName),
		"data":  newWidget,
	})
}

func handleCreateMenu(c *fiber.Ctx, db *gorm.DB, message string) error {
	parts := strings.Split(message, " ")
	if len(parts) < 3 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_COMMAND", "Please specify a menu name.")
	}
	menuName := strings.Join(parts[2:], " ")

	newMenu := models.Menuaccess{
		Menuname:     menuName,
		Menucode:     strings.ToLower(strings.ReplaceAll(menuName, " ", "_")),
		Description:  "Created via AI Assistant",
		Recordstatus: 1,
		Moduleid:     1, // Default
		Menutype:     "list",
		Sortorder:    99,
	}

	if err := db.Create(&newMenu).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error())
	}

	return helpers.SuccessResponse(c, "MENU_CREATED", fiber.Map{
		"reply": fmt.Sprintf("Success! Menu '%s' has been created.", menuName),
		"data":  newMenu,
	})
}

func handleCreateWorkflow(c *fiber.Ctx, db *gorm.DB, message string) error {
	parts := strings.Split(message, " ")
	if len(parts) < 3 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_COMMAND", "Please specify a workflow name.")
	}
	wfName := strings.Join(parts[2:], " ")

	newWf := models.Workflow{
		Wfname:       wfName,
		Wfdesc:       "Created via AI Assistant",
		Recordstatus: 1,
	}

	if err := db.Create(&newWf).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error())
	}

	return helpers.SuccessResponse(c, "WORKFLOW_CREATED", fiber.Map{
		"reply": fmt.Sprintf("Success! Workflow '%s' has been created.", wfName),
		"data":  newWf,
	})
}

func AiUploadHandler(c *fiber.Ctx, db *gorm.DB) error {
	fmt.Println("[AiUploadHandler] Started handling upload request")
	
	// Get requested component (default to 'upload')
	componentRequest := c.FormValue("component")
	if componentRequest == "" {
		componentRequest = "upload"
	}

	fmt.Printf("[AiUploadHandler] Using component: %s\n", componentRequest)

	var params []dbgenerator.WorkflowDetailResult

	if componentRequest == "document" {
		params = []dbgenerator.WorkflowDetailResult{
			{InputName: "action", CompValue: "upload_and_extract"},
			{InputName: "file_field", CompValue: "file"},
		}
	} else {
		// Generic upload defaults
		params = []dbgenerator.WorkflowDetailResult{
			{InputName: "filefield", CompValue: "file"},
			{InputName: "destination", CompValue: "public/chat_uploads"}, // Separate folder for chat uploads
			{InputName: "rename", CompValue: "true"},
		}
	}
	
	// Create a mock WorkflowContext
	ctx := &dbgenerator.WorkflowContext{
		FiberCtx: c,
		DB:       db,
		Params:   params,
		Extras:   make(map[string]interface{}),
	}

	// Get the component handler
	handler, ok := dbgenerator.GetComponent(componentRequest)
	if !ok {
		fmt.Printf("[AiUploadHandler] Error: Component '%s' not found in registry\n", componentRequest)
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "COMPONENT_ERROR", "Component not found: "+componentRequest)
	}

	fmt.Println("[AiUploadHandler] Component found, executing...")

	// Execute the component
	if err := handler.Execute(ctx); err != nil {
		fmt.Printf("[AiUploadHandler] Error executing component: %v\n", err)
		return err
	}

	fmt.Println("[AiUploadHandler] Execution successful")
	return nil
}
