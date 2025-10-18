// admin/auth.go
package admin

import (
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Language string `json:"language"`
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	auth := app.Group("/auth")

	// Public routes
	auth.Post("/login", func(c *fiber.Ctx) error { return LoginHandler(c, db) })

	// Protected routes
	auth.Use(AuthMiddleware)
	{
		//auth
		auth.Post("/logout", func(c *fiber.Ctx) error { return LogoutHandler(c, db) })
		auth.Get("/me", func(c *fiber.Ctx) error { return MeHandler(c, db) })
	}

	admin := app.Group("/admin")
	admin.Use(AuthMiddleware)
	{
		admin.Post("/generate-table", func(c *fiber.Ctx) error { return GenerateTableHandler(c, db) })
		admin.Post("/generate-multi-table", func(c *fiber.Ctx) error { return GenerateMultiTableHandler(c, db) })
		admin.Post("/generate-module", func(c *fiber.Ctx) error { return CreateModulesHandler(c, db) })
		plugin.RegisterModelRoutes(admin, db, models.Language{}, "language")
		plugin.RegisterModelRoutes(admin, db, models.Useraccess{}, "useraccess")
	}
}
