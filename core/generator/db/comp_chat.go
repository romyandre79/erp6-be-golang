package generator

import (
	"encoding/json"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterComponent("chat", func(ctx *WorkflowContext) error {
		return handleChat(ctx)
	})
}

func handleChat(ctx *WorkflowContext) error {
	var (
		action   string
		targetID int
		senderID int
		message  string
		attachment string
		msgID    int
	)

	// Extract parameters
	fmt.Printf("[CompChat] Params count: %d\n", len(ctx.Params))
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		
		// Fallback: If value is empty, try to resolve by InputName (e.g. $message)
		if val == "" {
			val = ResolveParam(ctx.FiberCtx, p.InputName)
			if strings.HasPrefix(val, "$") { val = "" } // ResolveParam returns "$name" if not found
		}

		fmt.Printf("[CompChat] Param '%s' resolved to: '%s'\n", p.InputName, val)

		switch strings.ToLower(p.InputName) {
		case "action":
			action = strings.ToLower(val)
		case "targetid", "target_id":
			targetID, _ = strconv.Atoi(val)
		case "senderid", "sender_id":
			senderID, _ = strconv.Atoi(val)
		case "message":
			message = val
		case "attachment":
			attachment = val
		case "messageid", "message_id":
			msgID, _ = strconv.Atoi(val)
		}
	}

	// Fallback/Defaults
	if senderID == 0 {
		if uid, ok := ctx.FiberCtx.Locals("userid").(int); ok {
			senderID = uid
		}
	}

	// Try to resolve action if empty (fallback to variable)
	// (Already handled in loop above, but double check)
	if action == "" {
		action = strings.ToLower(ResolveParam(ctx.FiberCtx, "action"))
		if strings.HasPrefix(action, "$") { action = "" }
	}

	fmt.Printf("[CompChat] Executing Action: '%s', Sender: %d, Target: %d, Message: '%s'\n", action, senderID, targetID, message)

	db := ctx.DB

	switch action {
	case "send":
		if senderID == 0 || targetID == 0 {
			fmt.Println("[CompChat] Error: SenderID or TargetID is 0")
			return fmt.Errorf("sender_id and target_id are required for send action")
		}
		if message == "" {
			fmt.Println("[CompChat] Error: Message is empty")
			return fmt.Errorf("message is required for send action")
		}

		// 1. Save to DB
		chatLog := models.Chat{
			SenderID:   senderID,
			ReceiverID: targetID,
			Message:    message,
			Attachment: attachment, // Added: attachment to DB insert
			IsRead:     0,
			CreatedAt:  time.Now(),
		}
		if err := db.Create(&chatLog).Error; err != nil {
			fmt.Printf("[CompChat] DB Create Error: %v\n", err)
			return err
		}
		fmt.Printf("[CompChat] Message saved to DB. ID: %d\n", chatLog.ChatID)

		// 2. Send via WebSocket
		if ws.GlobalHub != nil {
			// Payload for receiver matching AiAssistant.vue expectation
			payload, _ := json.Marshal(map[string]interface{}{
				"type":     "chat",
				"senderid": senderID,
				"data": map[string]interface{}{
					"text":       message,
					"attachment": attachment, // Added: attachment to WS payload
					"timestamp":  chatLog.CreatedAt,
				},
			})
			ws.GlobalHub.SendToUser(targetID, payload)
			fmt.Printf("[CompChat] WebSocket payload sent to User %d\n", targetID)
			
			// Optional: Echo back to sender? (Usually frontend handles optimistic UI, but good for confirmation)
			// ws.GlobalHub.SendToUser(senderID, payload) 
		} else {
			fmt.Println("[CompChat] Warning: ws.GlobalHub is nil")
		}

		// Result for workflow
		ctx.FiberCtx.Locals("wfResult", chatLog) // Legacy support?
		// Better result structure
		return nil

	case "gethistory":
		if senderID == 0 || targetID == 0 {
			return fmt.Errorf("sender_id and target_id are required for gethistory action")
		}

		var chats []models.Chat
		err := db.Where(
			db.Where("senderid = ? AND receiverid = ?", senderID, targetID).
				Or("senderid = ? AND receiverid = ?", targetID, senderID),
		).Order("createdat ASC").Find(&chats).Error

		if err != nil {
			return err
		}

		// Store result in context so next node can use it (or End node returns it)
		// Usually InternalFlow puts result in wfEngine. We just return it via some mechanism?
		// The ComponentHandler signature only returns error. 
		// But in executeflow.go, it looks at where?
		// Ah, standard components don't easily "return" data unless they write to Locals or specific structs.
		// Wait, external runner wrote to wfEngine.
		// Let's see how `comp_decision` or others do it.
		// `comp_table` writes to `c.Locals("query_result", ...)` or returns map.
		
		// Wait, the interface is `Execute(ctx) error`.
		// But `InternalFlow` uses the return? No.
		// `InternalFlow` captures `err`.
		
		// Let's look at `external_runner.go` again.
		// It manually appends to `wfEngine` in Locals.
		
		wm := WorkflowEngine{
			ResultNode:    chats,
		}
		
		wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)

		return nil

	case "markread":
		if msgID == 0 {
			// Maybe mark all read for a user?
			if senderID != 0 && targetID != 0 {
				// senderID here is the "viewer" (me), targetID is the "other person"
				// Update all messages FROM targetID TO me as read
				db.Model(&models.Chat{}).Where("senderid = ? AND receiverid = ?", targetID, senderID).Update("isread", 1)
				return nil
			}
			return fmt.Errorf("messageid is required for markread action")
		}
		// Update specific message
		db.Model(&models.Chat{}).Where("chatid = ?", msgID).Update("isread", 1)
		return nil

	case "getuserlist":
		// Get list of users except me
		var users []models.Useraccess
		if err := db.Where("useraccessid != ? AND recordstatus = 1", senderID).Order("realname ASC").Find(&users).Error; err != nil {
			return err
		}
		
		// Return specific fields
		var result []map[string]interface{}
		for _, u := range users {
			photo := u.Userphoto
			if photo == "" { photo = "" } // Ensure string
			
			result = append(result, map[string]interface{}{
				"useraccessid": u.Useraccessid,
				"username":     u.Username,
				"realname":     u.Realname,
				"email":        u.Email,
				"userphoto":    photo,
				"isonline":     u.Isonline,
			})
		}
		
		wm := WorkflowEngine{
			ResultNode:    result,
		}
		wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
		
		return nil

	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}
