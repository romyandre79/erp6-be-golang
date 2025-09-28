// plugins/hello/hello.go
package admin

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
)

type AdminPlugin struct{}

func (h AdminPlugin) Name() string { return "admin" }

func (h AdminPlugin) Init() error { return nil }

func (h AdminPlugin) Routes(app *fiber.App) {
	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Admin plugin!",
		})
	})
}

func init() {
	plugin.Register(AdminPlugin{})
}
