package admin

import (
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/scheduler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ReloadSchedulerHandler manually triggers scheduler sync from workflows
func ReloadSchedulerHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")

	// Check permission (optional - you can remove this if you want it accessible to all)
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	// Trigger scheduler sync
	scheduler.SyncJobsFromWorkflows(db)
	scheduler.LoadJobs(db)

	return helpers.SuccessResponse(c, "SCHEDULER_RELOADED", fiber.Map{
		"message": "Scheduler has been reloaded from workflows",
	})
}
