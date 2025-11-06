package admin

import (
	"erp6-be-golang/core/events"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitWidget() {
	events.Register("BeforeGet:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:widget", func(data interface{}) error {
		return nil
	})
}

func DashboardListHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")
	moduleName := c.Query("module")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	var menus []response.Widget
	err := db.
		Table("widget a").
		Select(`
			DISTINCT a.widgetid, a.widgetname, a.widgettitle, a.widgetversion, a.widgetform, a.widgetby, a.description, dashgroup, position, a.moduleid, modulename
		`).
		Joins("JOIN userdash b ON b.widgetid = a.widgetid").
		Joins("JOIN usergroup c ON c.groupaccessid = b.groupaccessid").
		Joins("JOIN modules d ON d.moduleid = a.moduleid").
		Where("c.useraccessid = ? AND d.modulename = ? AND a.recordstatus = 1", userID, moduleName).
		Order("b.dashgroup, b.position").
		Scan(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", menus)
}
