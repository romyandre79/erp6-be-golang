// admin/auth.go
package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	auth := app.Group("/auth")

	// Public routes
	auth.Post("/login", func(c *fiber.Ctx) error { return loginHandler(c, db) })

	// Protected routes
	auth.Use(AuthMiddleware)
	{
		//auth
		auth.Post("/logout", func(c *fiber.Ctx) error { return logoutHandler(c, db) })
		auth.Get("/me", func(c *fiber.Ctx) error { return meHandler(c, db) })
	}

	admin := app.Group("/admin")
	admin.Use(AuthMiddleware)
	{
		admin.Post("/generate-table", func(c *fiber.Ctx) error { return generateTableHandler(c, db) })
		admin.Post("/generate-module", func(c *fiber.Ctx) error { return createModulesHandler(c, db) })
	}
}
