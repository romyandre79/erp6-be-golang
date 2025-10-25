package blog

import (
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/models"
	"erp6-be-golang/plugins/admin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes digunakan untuk mendaftarkan route plugin blog.
// Isi sesuai kebutuhan modul.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// TODO: implement routes
	blogAuth := app.Group("/blogAuth")
	blogAuth.Use(admin.AuthMiddleware)
	{
		plugin.RegisterModelRoutes(blogAuth, db, models.Category{}, "category")
		plugin.RegisterModelRoutes(blogAuth, db, models.Post{}, "post")
		plugin.RegisterModelRoutes(blogAuth, db, models.Postcomment{}, "postcomment")
	}
	blog := app.Group("/blog")
	blog.Get("/post", func(c *fiber.Ctx) error { return BlogPostHandler(c, db) })
	blog.Get("/post/slug", func(c *fiber.Ctx) error { return BlogSlugPostHandler(c, db) })

}
