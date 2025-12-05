package generator

import (
	"encoding/json"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("sendmessage", func(ctx *WorkflowContext) error {
		return handleSendMessage(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

// handleSendMessage sends a notification to a user via WebSocket and saves it to the database
// Supported parameters:
// - sendto: User ID to send the message to (required)
// - message: The message content (required)
// - title: The title of the notification (optional)
// - docno: Document number associated with the notification (optional)
func handleSendMessage(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		sendTo  int
		message string
		title   = "Notification"
		docNo   string
	)

	// Extract parameters from workflow
	for _, p := range params {
		val := strings.TrimSpace(p.CompValue)
		switch strings.ToLower(p.InputName) {
		case "sendtonotif":
			// Try to parse as int
			if v, err := strconv.Atoi(val); err == nil {
				sendTo = v
			} else {
				// If strictly required to be int
				return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "sendto must be a numeric user ID")
			}
		case "messagenotif":
			message = val
		case "titlenotif":
			if val != "" {
				title = val
			}
		case "othernotif":
			docNo = val
		}
	}

	// Validate required parameters
	if sendTo == 0 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "sendto is required")
	}
	if message == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "message is required")
	}

	// Create UserTodo record
	todo := models.Usertodo{
		Useraccessid: sendTo,
		Menuname:     title,
		Description:  message,
		Docno:        docNo,
		Isread:       0,
		// Tododate might be auto-set by DB or GORM hook? Usually it is.
		// If not, we might need to set it. But notification.go didn't set it.
	}

	if err := db.Create(&todo).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to save notification: %v", err))
	}

	// Broadcast via WS if GlobalHub is available
	if ws.GlobalHub != nil {
		payload, err := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": todo,
		})
		if err == nil {
			ws.GlobalHub.SendToUser(sendTo, payload)
		} else {
			// Log error?
			fmt.Printf("Error marshaling notification payload: %v\n", err)
		}
	} else {
		fmt.Println("Warning: GlobalHub is nil, cannot send WebSocket notification")
	}

	return nil
}
