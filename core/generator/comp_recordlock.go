package generator

import (
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("lockrecord", func(ctx *WorkflowContext) error {
		return handleLockRecord(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
	RegisterComponent("unlockrecord", func(ctx *WorkflowContext) error {
		return handleUnlockRecord(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

// handleLockRecord locks a record for editing
// Parameters:
// - tablename: Name of the table (required)
// - recordid: ID of the record to lock (required)
// - locktype: Type of lock (optional, default: "edit")
func handleLockRecord(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		tableName string
		recordID  int
		lockType  = "edit"
	)

	// Extract parameters
	for _, p := range params {
		val := strings.TrimSpace(ResolveParam(c, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "tablename", "table":
			tableName = val
		case "recordid", "record_id", "id":
			if v, err := strconv.Atoi(val); err == nil {
				recordID = v
			}
		case "locktype", "lock_type", "type":
			if val != "" {
				lockType = val
			}
		}
	}

	// Validate required parameters
	if tableName == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "tablename is required")
	}
	if recordID == 0 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "recordid is required")
	}

	// Get current user ID
	userID, ok := c.Locals("userid").(int)
	if !ok || userID == 0 {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	// Get session ID
	sessionID := c.Get("X-Session-ID", "")

	// Check if record is already locked
	var existingLock models.Recordlock
	err := db.Where("tablename = ? AND recordid = ?", tableName, recordID).First(&existingLock).Error
	
	if err == nil {
		// Record is already locked
		if existingLock.Lockedby == userID {
			// Same user, refresh the lock
			existingLock.Lockedat = time.Now()
			existingLock.Sessionid = sessionID
			if err := db.Save(&existingLock).Error; err != nil {
				return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to refresh lock: %v", err))
			}
			fmt.Printf("[LockRecord] Lock refreshed: %s#%d by user %d\n", tableName, recordID, userID)
			return nil
		}

		// Locked by another user - get user info
		var lockingUser models.Useraccess
		db.Where("useraccessid = ?", existingLock.Lockedby).First(&lockingUser)

		return helpers.FailResponse(c, fiber.StatusConflict, "RECORD_LOCKED", 
			fmt.Sprintf("Record is being edited by %s since %s", 
				lockingUser.Username, existingLock.Lockedat.Format("15:04:05")))
	}

	// Create new lock
	lock := models.Recordlock{
		Tablename: tableName,
		Recordid:  recordID,
		Lockedby:  userID,
		Locktype:  lockType,
		Sessionid: sessionID,
	}

	if err := db.Create(&lock).Error; err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to create lock: %v", err))
	}

	fmt.Printf("[LockRecord] Record locked: %s#%d by user %d\n", tableName, recordID, userID)

	// Broadcast lock notification to all users
	if ws.GlobalHub != nil {
		// Get all users from database
		var users []models.Useraccess
		db.Find(&users)

		// Get current user info
		var currentUser models.Useraccess
		db.Where("useraccessid = ?", userID).First(&currentUser)

		// Send notification to all users except the one who locked
		for _, user := range users {
			if user.Useraccessid != userID {
				// Create notification in usertodo
				todo := models.Usertodo{
					Useraccessid: user.Useraccessid,
					Menuname:     "Record Locked",
					Description:  fmt.Sprintf("%s is editing %s #%d", currentUser.Username, tableName, recordID),
					Isread:       0,
				}
				db.Create(&todo)

				// Send WebSocket notification
				ws.GlobalHub.SendToUser(user.Useraccessid, []byte(fmt.Sprintf(`{
					"type": "record_locked",
					"tablename": "%s",
					"recordid": %d,
					"locked_by": %d,
					"locked_by_name": "%s",
					"locked_at": "%s"
				}`, tableName, recordID, userID, currentUser.Username, time.Now().Format(time.RFC3339))))
			}
		}
		fmt.Printf("[LockRecord] Lock notification broadcast to all users\n")
	}

	return nil
}

// handleUnlockRecord releases a lock on a record
// Parameters:
// - tablename: Name of the table (required)
// - recordid: ID of the record to unlock (required)
func handleUnlockRecord(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		tableName string
		recordID  int
	)

	// Extract parameters
	for _, p := range params {
		val := strings.TrimSpace(ResolveParam(c, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "tablename", "table":
			tableName = val
		case "recordid", "record_id", "id":
			if v, err := strconv.Atoi(val); err == nil {
				recordID = v
			}
		}
	}

	// Validate required parameters
	if tableName == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "tablename is required")
	}
	if recordID == 0 {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "recordid is required")
	}

	// Get current user ID
	userID, ok := c.Locals("userid").(int)
	if !ok || userID == 0 {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	// Delete lock (only if locked by current user)
	result := db.Where("tablename = ? AND recordid = ? AND lockedby = ?", tableName, recordID, userID).
		Delete(&models.Recordlock{})

	if result.Error != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("failed to unlock record: %v", result.Error))
	}

	if result.RowsAffected == 0 {
		fmt.Printf("[UnlockRecord] No lock found for %s#%d by user %d\n", tableName, recordID, userID)
		return nil // Not an error, just no lock to release
	}

	fmt.Printf("[UnlockRecord] Record unlocked: %s#%d by user %d\n", tableName, recordID, userID)

	// Broadcast unlock notification to all users
	if ws.GlobalHub != nil {
		var users []models.Useraccess
		db.Find(&users)

		var currentUser models.Useraccess
		db.Where("useraccessid = ?", userID).First(&currentUser)

		for _, user := range users {
			if user.Useraccessid != userID {
				ws.GlobalHub.SendToUser(user.Useraccessid, []byte(fmt.Sprintf(`{
					"type": "record_unlocked",
					"tablename": "%s",
					"recordid": %d,
					"unlocked_by": %d,
					"unlocked_by_name": "%s"
				}`, tableName, recordID, userID, currentUser.Username)))
			}
		}
		fmt.Printf("[UnlockRecord] Unlock notification broadcast to all users\n")
	}

	return nil
}
