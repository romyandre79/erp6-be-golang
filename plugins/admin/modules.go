package admin

import (
	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/core/helpers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func createModulesHandler(c *fiber.Ctx, db *gorm.DB) error {
	return helpers.SuccessResponse(c, "DATA_SAVED", nil)
}

func generateTableHandler(c *fiber.Ctx, db *gorm.DB) error {
	menuName := c.FormValue("menu")
	tableName := c.FormValue("table")
	IsPermission, err := CheckUserPermission(c, db, menuName, PermWrite)
	if err != nil || !IsPermission {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_AUTHORIZE", err.Error())
	}

	if tableName == "" {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_TABLE", "")
	}

	err = generator.GenerateStructWithMeta(db, tableName)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_GENERAL_TABLE", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA_SAVED", nil)
}
