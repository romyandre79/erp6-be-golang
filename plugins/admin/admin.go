package admin

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
)

type AdminPlugin struct{}

func (h AdminPlugin) Name() string { return "admin" }

func (h AdminPlugin) Init() error { return nil }

func (h AdminPlugin) Routes(app *fiber.App) {
	RegisterRoutes(app)
}

func init() {
	plugin.Register(AdminPlugin{})
}
