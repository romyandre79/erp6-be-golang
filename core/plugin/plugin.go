// core/plugin/plugin.go
package plugin

import "github.com/gofiber/fiber/v2"

type Plugin interface {
	Name() string
	Init() error
	Routes(app *fiber.App)
}
