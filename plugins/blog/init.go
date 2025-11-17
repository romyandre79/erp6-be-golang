package blog

import (
	"erp6-be-golang/core/plugin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BlogPlugin struct{}

func (h BlogPlugin) Name() string { return "blog" }

func (h BlogPlugin) Init(db *gorm.DB) error { return nil }

func (h BlogPlugin) Routes(app *fiber.App, db *gorm.DB) {
	RegisterRoutes(app, db)
}

func init() {
	plugin.Register(BlogPlugin{})
}
