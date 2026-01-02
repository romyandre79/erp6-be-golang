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
		messageType = "notification" // Default to notification
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
		var wfEngine []WorkflowEngine
		var ok bool
		
		if c != nil {
			wfEngine, ok = c.Locals("wfEngine").([]WorkflowEngine)
		} else {
			// WhatsApp context - check ctx.Extras if available
			// Note: ctx is not available here, so we can't access Extras
			// This will be handled later
		}
		
		if ok && len(wfEngine) > 0 {
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
		var wfEngine []WorkflowEngine
		var ok bool
		
		if c != nil {
			wfEngine, ok = c.Locals("wfEngine").([]WorkflowEngine)
		}
		
		if ok && len(wfEngine) > 0 {
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

	// If sendTo is still 0, look for authenticated user in Locals
	if sendTo == 0 {
		if uid, ok := c.Locals("userid").(int); ok && uid > 0 {
			sendTo = uid
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
	
	// PRIORITY: Check for WhatsApp/Telegram callback in Locals (via comp_internal_*.go)
	// This ensures messages sourced from WA/TG are routed back to them immediately,
	// regardless of messageType (chat vs notification), preventing leaks to Web/DB.
	if c != nil {
		// Method 1: Direct Local (Telegram/New)
		if callback, ok := c.Locals("send_wa_callback").(func(string)); ok {
			fmt.Printf("[SendMessage] triggering WA/TG callback (Direct): %s\n", message)
			callback(message)
			return nil
		}
		
		// Method 2: wfExtras Map (WhatsApp/Legacy)
		if wfExtras, ok := c.Locals("wfExtras").(map[string]interface{}); ok {
			if callback, ok := wfExtras["send_wa_callback"].(func(string)); ok {
				fmt.Printf("[SendMessage] triggering WA/TG callback (Extras): %s\n", message)
				callback(message)
				return nil // Stop execution here
			}
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
		// Check if this is a WhatsApp context (nil FiberCtx)
		if c == nil {
			// WhatsApp context - store message in wfEngine for WhatsApp handler to send
			fmt.Printf("[SendMessage] WhatsApp context detected, storing message: %s\n", message)
			return nil
		}
		
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
		
		// Modify message if workflow is executing AND we don't have actual data yet
		displayMessage := message
		if executeFlag == "true" && (message == "" || message == "Data Customer sent" || strings.Contains(message, "processed successfully")) {
			// Only show processing message if we don't have actual data
			displayMessage = "⏳ Processing your request, please wait..."
		}
		
		if ws.GlobalHub != nil {
			payload, err := json.Marshal(map[string]interface{}{
				"type":               "chat",
				"senderid":           sendTo, // User who sent the message
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

	// Handle Telegram Message
	// Check if we have source=telegram and chat_id in the workflow engine history
	var tgChatID int64
	var tgSource bool
	
	// Scan previous nodes for Telegram metadata
	if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok {
		for i := len(wfEngine) - 1; i >= 0; i-- {
			if wfEngine[i].ResultNode != nil {
				if resultMap, ok := wfEngine[i].ResultNode.(map[string]interface{}); ok {
					// Check source
					if s, ok := resultMap["source"].(string); ok && s == "telegram" {
						tgSource = true
					}
					// Check chat_id (might be int, int64, or float64 from JSON)
					if cid, ok := resultMap["chat_id"]; ok {
						if v, ok := cid.(int64); ok {
							tgChatID = v
						} else if v, ok := cid.(int); ok {
							tgChatID = int64(v)
						} else if v, ok := cid.(float64); ok {
							tgChatID = int64(v)
						}
					}
				}
			}
		}
	}

	if tgSource && tgChatID != 0 {
		fmt.Printf("[SendMessage] Sending to Telegram ChatID: %d\n", tgChatID)
		if err := SendTelegramMessage(tgChatID, message); err != nil {
			fmt.Printf("[SendMessage] Telegram Error: %v\n", err)
			return err
		}
		
		// Save conversation state if available (Required for AI flow)
		if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok && len(wfEngine) > 0 {
			lastResult := wfEngine[len(wfEngine)-1]
			if lastResult.ResultNode != nil {
				if resultMap, ok := lastResult.ResultNode.(map[string]interface{}); ok {
					if conversationState, ok := resultMap["conversation_state"].(string); ok {
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
								fmt.Printf("[SendMessage] Saved Telegram conversation state for user %d\n", userID)
							}
						}
					}
				}
			}
		}
		
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
		Tododate:     time.Now(), // Explicitly set time to ensure it's in the JSON
	}

	if err := db.Create(&todo).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to save notification: %v", err))
	}
	fmt.Printf("[SendMessage] Notification saved to DB: ID=%d\n", todo.Usertodoid)

	// Broadcast via WS if GlobalHub is available
	if ws.GlobalHub != nil {
		payloadMap := map[string]interface{}{
			"type": "notification",
			"data": todo,
		}
		
		payload, err := json.Marshal(payloadMap)
		if err == nil {
			ws.GlobalHub.SendToUser(sendTo, payload)
			fmt.Printf("[SendMessage] Notification sent via WebSocket to user %d. Payload: %s\n", sendTo, string(payload))
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
