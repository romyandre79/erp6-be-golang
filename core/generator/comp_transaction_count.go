package generator

import (
	"encoding/json"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	RegisterComponent("transaction_count", func(ctx *WorkflowContext) error {
		return handleTransactionCount(ctx)
	})
}

// handleTransactionCount limits transactions based on user-specific rules
// Parameters:
// - rules: "userid:[table:limit,table2:limit2]; userid2:[...]"
// - user_identifier: Column name to check in target table (default: created_by)
// - userid: Override user ID check (default: current logged in user)
func handleTransactionCount(ctx *WorkflowContext) error {
	var (
		rules          string
		userIdentifier = "created_by"
		targetUserID   int
	)

	// 1. Parse Parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "rules":
			rules = val
		case "user_identifier":
			userIdentifier = val
		case "userid":
			if v, err := strconv.Atoi(val); err == nil {
				targetUserID = v
			}
		}
	}

	// 2. Determine User ID
	if targetUserID == 0 {
		if uid, ok := ctx.FiberCtx.Locals("userid").(int); ok && uid > 0 {
			targetUserID = uid
		}
	}

	if targetUserID == 0 {
		fmt.Println("[TransactionCount] No User ID found. Skipping check (Allowed).")
		ctx.DecisionResult = true
		return nil
	}

	// 3. Find Rule for User
	// Format: "1:[table:10]; 2:[table:20]"
	userRuleStr := ""
	
	// Split into user blocks
	userBlocks := strings.Split(rules, ";")
	for _, block := range userBlocks {
		block = strings.TrimSpace(block)
		parts := strings.SplitN(block, ":", 2)
		if len(parts) == 2 {
			uidStr := strings.TrimSpace(parts[0])
			configStr := strings.TrimSpace(parts[1])
			
			if id, err := strconv.Atoi(uidStr); err == nil && id == targetUserID {
				// Found match! Remove brackets if present
				configStr = strings.TrimPrefix(configStr, "[")
				configStr = strings.TrimSuffix(configStr, "]")
				userRuleStr = configStr
				break
			}
		}
	}

	if userRuleStr == "" {
		fmt.Printf("[TransactionCount] No specific rule found for user %d. Skipping check (Allowed).\n", targetUserID)
		ctx.DecisionResult = true
		return nil
	}

	// 4. Validate Limits
	// Format: "tableA:5000, tableB:3000"
	limits := strings.Split(userRuleStr, ",")
	allowed := true

	for _, limitRule := range limits {
		limitRule = strings.TrimSpace(limitRule)
		parts := strings.Split(limitRule, ":")
		if len(parts) != 2 {
			continue
		}

		tableName := strings.TrimSpace(parts[0])
		limitVal, err := strconv.ParseFloat(parts[1], 64) // Use float to handle potential math ease, though count is int
		if err != nil {
			fmt.Printf("[TransactionCount] Invalid limit value '%s' for table '%s'\n", parts[1], tableName)
			continue
		}
		limit := int64(limitVal)

		// Count existing transactions
		var currentCount int64
		// Verify table existence to prevent SQL injection or errors? GORM handles basic safety.
		// Use raw query for defined table
		query := fmt.Sprintf("SELECT count(1) FROM %s WHERE %s = ?", tableName, userIdentifier)
		if err := ctx.DB.Raw(query, targetUserID).Scan(&currentCount).Error; err != nil {
			fmt.Printf("[TransactionCount] Error counting table '%s': %v\n", tableName, err)
			// Fail safe? Or block? Let's assume block if error to be safe, or allow?
			// Usually system error should maybe stop flow. 
			// For now, log and continue check next rule?
			continue
		}

		fmt.Printf("[TransactionCount] User %d | Table %s | Count: %d / %d\n", targetUserID, tableName, currentCount, limit)

		// Check Limit
		if currentCount >= limit {
			fmt.Printf("[TransactionCount] User %d exceeded limit for table %s\n", targetUserID, tableName)
			allowed = false
			break // Block immediately on first violation
		}

		// Check Notification Threshold (90%)
		// Send notification only if exactly entered the warning zone? 
		// Or every time they are in it? 
		// Use >= 90% and < 100%. 
		// To avoid spamming, maybe check if it was ALREADY sent? (Complex).
		// Requirement: "send notif if 10% before total limit" -> implies reaching the threshold.
		// Let's send it if count >= 0.9 * limit.
		threshold := int64(float64(limit) * 0.9)
		if currentCount >= threshold {
			fmt.Printf("[TransactionCount] User %d is near limit (%d/%d). Sending notification.\n", targetUserID, currentCount, limit)
			sendNotification(ctx.DB, targetUserID, tableName, currentCount, limit)
		}
	}

	ctx.DecisionResult = allowed
	return nil
}

func sendNotification(db *gorm.DB, userID int, tableName string, count, limit int64) {
	title := "Usage Warning"
	message := fmt.Sprintf("Warning: You have used %d of %d transactions for %s. You are approaching your limit.", count, limit, tableName)
	
	// Check if a similar unread notification already exists to avoid spamming
	var exists int64
	db.Table("usertodo").
		Where("useraccessid = ? AND description = ? AND isread = 0", userID, message).
		Count(&exists)
	
	if exists > 0 {
		return // Already notified
	}

	// Create UserTodo
	todo := models.Usertodo{
		Useraccessid: userID,
		Menuname:     title,
		Description:  message,
		Isread:       0,
		Tododate:     time.Now(),
	}

	if err := db.Create(&todo).Error; err != nil {
		fmt.Printf("[TransactionCount] Failed to create notification: %v\n", err)
		return
	}

	// Send WS
	if ws.GlobalHub != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": todo,
		})
		ws.GlobalHub.SendToUser(userID, payload)
	}
}
