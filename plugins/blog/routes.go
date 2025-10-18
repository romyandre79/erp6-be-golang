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
	blog := app.Group("/blog")
blog.Use(admin.AuthMiddleware)
{
plugin.RegisterModelRoutes(blog, db, models.Category{}, "category")
plugin.RegisterModelRoutes(blog, db, models.Post{}, "post")
plugin.RegisterModelRoutes(blog, db, models.Postcomment{}, "postcomment")

}

}
