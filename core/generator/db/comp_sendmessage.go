package generator

import (
	"encoding/json"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
		sendTo      int
		message     string
		title       = "Notification"
		docNo       string
		messageType = "notification" // Default to notification, can be "chat" or "notification"
	)

	// Extract parameters from workflow
	for _, p := range params {
		val := strings.TrimSpace(ResolveParam(c, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "sendtonotif", "user_id":
			// Try to parse as int
			if v, err := strconv.Atoi(val); err == nil {
				sendTo = v
			}
			// Don't return error here - allow fallback to form data or node result
		case "messagenotif", "message":
			message = val
		case "titlenotif", "title":
			if val != "" {
				title = val
			}
		case "othernotif", "docno":
			docNo = val
		case "message_type", "delivery_type", "messagetype":
			if val != "" {
				messageType = strings.ToLower(val)
			}
		}
	}

	// If user_id not found in workflow params, try to get it from form data
	if sendTo == 0 {
		if userIDStr := c.FormValue("user_id"); userIDStr != "" {
			if v, err := strconv.Atoi(userIDStr); err == nil {
				sendTo = v
			}
		}
	}

	// If message is still empty, try to get it from previous node's result
	if message == "" {
		if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok && len(wfEngine) > 0 {
			// Get the last result from previous node
			lastResult := wfEngine[len(wfEngine)-1]
			if lastResult.ResultNode != nil {
				// Try to extract 'message' field from result
				if resultMap, ok := lastResult.ResultNode.(map[string]interface{}); ok {
					if msg, ok := resultMap["message"].(string); ok && msg != "" {
						message = msg
					}
				}
			}
		}
	}

	// If user_id is still 0, try to get it from previous node's result
	if sendTo == 0 {
		if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok && len(wfEngine) > 0 {
			// Get the last result from previous node
			lastResult := wfEngine[len(wfEngine)-1]
			if lastResult.ResultNode != nil {
				// Try to extract 'user_id' field from result
				if resultMap, ok := lastResult.ResultNode.(map[string]interface{}); ok {
					if userIDStr, ok := resultMap["user_id"].(string); ok && userIDStr != "" {
						if v, err := strconv.Atoi(userIDStr); err == nil {
							sendTo = v
						}
					}
				}
			}
		}
	}

	// Debug logging
	fmt.Printf("[SendMessage] Extracted values: sendTo=%d, message='%s', title='%s', type='%s'\n", sendTo, message, title, messageType)

	// Replace template variables in message
	if commandOutput := c.Locals("commandOutput"); commandOutput != nil {
		if outputStr, ok := commandOutput.(string); ok {
			message = strings.ReplaceAll(message, "{{commandOutput}}", outputStr)
			title = strings.ReplaceAll(title, "{{commandOutput}}", outputStr)
			fmt.Printf("[SendMessage] Replaced {{commandOutput}} template variable\n")
		}
	}

	// Validate required parameters
	if message == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "message is required")
	}

	// If sendTo is 0, broadcast to all users
	if sendTo == 0 {
		fmt.Printf("[SendMessage] Broadcasting to all users: message='%s'\n", message)
		
		// Get all users from database
		var users []models.Useraccess
		if err := db.Find(&users).Error; err != nil {
			return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to get users: %v", err))
		}

		// Send to each user
		for _, user := range users {
			// Create notification in usertodo
			todo := models.Usertodo{
				Useraccessid: user.Useraccessid,
				Menuname:     title,
				Description:  message,
				Docno:        docNo,
				Isread:       0,
			}
			db.Create(&todo)

			// Send via WebSocket
			if ws.GlobalHub != nil {
				payload, _ := json.Marshal(map[string]interface{}{
					"type": messageType,
					"data": todo,
				})
				ws.GlobalHub.SendToUser(user.Useraccessid, payload)
			}
		}

		fmt.Printf("[SendMessage] Broadcast complete: sent to %d users\n", len(users))
		return nil
	}

	// Handle based on message type
	if messageType == "chat" {
		// For chat messages, just send via WebSocket without saving to DB
		// Also extract conversation_state and execute flag from previous node if available
		var conversationState string
		var executeFlag string
		if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok && len(wfEngine) > 0 {
			lastResult := wfEngine[len(wfEngine)-1]
			if lastResult.ResultNode != nil {
				if resultMap, ok := lastResult.ResultNode.(map[string]interface{}); ok {
					if state, ok := resultMap["conversation_state"].(string); ok {
						conversationState = state
						
						// Save conversation state to file for next message
						if userID, ok := c.Locals("userid").(int); ok && userID > 0 {
							conversationDir := "./tmp/ai_conversations"
							os.MkdirAll(conversationDir, 0755)
							conversationFile := fmt.Sprintf("%s/%d.json", conversationDir, userID)
							stateData := map[string]interface{}{
								"conversation_state": conversationState,
								"updated_at":         time.Now().Format(time.RFC3339),
							}
							if data, err := json.Marshal(stateData); err == nil {
								os.WriteFile(conversationFile, data, 0644)
								fmt.Printf("[SendMessage] Saved conversation state for user %d\n", userID)
							}
						}
					}
					
					// Check execute flag
					if exec, ok := resultMap["execute"].(string); ok {
						executeFlag = exec
					}
				}
			}
		}
		
		// Modify message if workflow is executing
		displayMessage := message
		if executeFlag == "true" {
			displayMessage = "⏳ Processing your request, please wait..."
		}
		
		if ws.GlobalHub != nil {
			payload, err := json.Marshal(map[string]interface{}{
				"type":               "chat",
				"message":            displayMessage,
				"title":              title,
				"conversation_state": conversationState,
				"executing":          executeFlag == "true",
			})
			if err == nil {
				ws.GlobalHub.SendToUser(sendTo, payload)
				fmt.Printf("[SendMessage] Chat message sent via WebSocket to user %d (executing=%s)\n", sendTo, executeFlag)
			} else {
				fmt.Printf("Error marshaling chat payload: %v\n", err)
			}
		} else {
			fmt.Println("Warning: GlobalHub is nil, cannot send WebSocket chat message")
		}
		fmt.Printf("[SendMessage] Success! Chat message sent to user %d\n", sendTo)
		return nil
	}

	// For notification type, save to DB and send via WebSocket
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
	fmt.Printf("[SendMessage] Notification saved to DB: ID=%d\n", todo.Usertodoid)

	// Broadcast via WS if GlobalHub is available
	if ws.GlobalHub != nil {
		payload, err := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": todo,
		})
		if err == nil {
			ws.GlobalHub.SendToUser(sendTo, payload)
			fmt.Printf("[SendMessage] Notification sent via WebSocket to user %d\n", sendTo)
		} else {
			// Log error?
			fmt.Printf("Error marshaling notification payload: %v\n", err)
		}
	} else {
		fmt.Println("Warning: GlobalHub is nil, cannot send WebSocket notification")
	}

	fmt.Printf("[SendMessage] Success! Notification sent to user %d\n", sendTo)
	return nil
}
