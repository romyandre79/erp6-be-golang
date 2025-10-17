package common

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommonPlugin struct{}

func (h CommonPlugin) Name() string { return "common" }

func (h CommonPlugin) Init(db *gorm.DB) error { return nil }

func (h CommonPlugin) Routes(app *fiber.App, db *gorm.DB) {
	RegisterRoutes(app, db)
}

func init() {
	plugin.Register(CommonPlugin{})
}
