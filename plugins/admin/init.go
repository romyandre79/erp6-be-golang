package admin

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminPlugin struct{}

func (h AdminPlugin) Name() string { return "admin" }

func (h AdminPlugin) Init(db *gorm.DB) error { return nil }

func (h AdminPlugin) Routes(app *fiber.App, db *gorm.DB) {
	RegisterRoutes(app, db)
}

func init() {
	plugin.Register(AdminPlugin{})
}
