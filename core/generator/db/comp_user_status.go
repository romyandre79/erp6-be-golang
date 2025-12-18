package generator

import (
	"erp6-be-golang/core/ws"
	"fmt"
)

func init() {
	fmt.Println("[Init] Registering component 'check_online' & aliases")
	// Register multiple aliases to be safe against casing in Flow Designer
	handler := func(ctx *WorkflowContext) error {
		return handleCheckOnline(ctx)
	}

	RegisterComponent("check_online", handler)
	RegisterComponent("checkonline", handler) 
	RegisterComponent("CheckOnline", handler) // Match exact case from error log
}

// handleCheckOnline synchronizes the 'isonline' status in the database
// with the actual active WebSocket connections.
func handleCheckOnline(ctx *WorkflowContext) error {
	fmt.Println("[CheckOnline] Handler START")

	// 1. Get list of actually connected user IDs from the Hub
	onlineUserIDs := ws.GlobalHub.GetOnlineUsers()
	
	countActive := len(onlineUserIDs)
	fmt.Printf("[CheckOnline] Found %d active WebSocket connections\n", countActive)

	// 2. Prepare query to set isonline = 0 for anyone NOT in this list
	// but currently marked as isonline = 1
	
	db := ctx.DB
	
	var result *WorkflowContext // Just to hold DB result if needed, but we use db directly
	_ = result

	if countActive == 0 {
		// If no one is online, set everyone to offline
		res := db.Table("useraccess").Where("isonline = 1").Update("isonline", 0)
		if res.Error != nil {
			fmt.Printf("[CheckOnline] Error updating all users to offline: %v\n", res.Error)
			return res.Error
		}
		fmt.Printf("[CheckOnline] No active users. Set %d users to offline.\n", res.RowsAffected)
	} else {
		// Set offline for users who are marked online BUT NOT in the active list
		res := db.Table("useraccess").
			Where("isonline = 1 AND useraccessid NOT IN ?", onlineUserIDs).
			Update("isonline", 0)
			
		if res.Error != nil {
			fmt.Printf("[CheckOnline] Error syncing online users: %v\n", res.Error)
			return res.Error
		}
		fmt.Printf("[CheckOnline] Sync complete. Set %d ghost users to offline.\n", res.RowsAffected)
	}

	// Optional: Return stats to the flow
	ctx.DecisionResult = true
	
	return nil
}
