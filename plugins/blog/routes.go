package blog

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes digunakan untuk mendaftarkan route plugin blog.
// Isi sesuai kebutuhan modul.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// TODO: implement routes
	blog := app.Group("/blog")
	blog.Get("/post", func(c *fiber.Ctx) error { return BlogPostHandler(c, db) })
	blog.Get("/post/slug", func(c *fiber.Ctx) error { return BlogSlugPostHandler(c, db) })
}
