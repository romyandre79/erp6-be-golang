// core/plugin/plugin.go
package plugin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Plugin interface {
	Name() string
	Init(db *gorm.DB) error
	Routes(app *fiber.App, db *gorm.DB)
}
