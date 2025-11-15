// admin/auth.go
package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var mediaRoot = "./public"

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
		auth.Get("/me", func(c *fiber.Ctx) error { return MeHander(c, db) })
	}

	admin := app.Group("/admin")
	admin.Use(AuthMiddleware)
	{
		admin.Get("/getmenu", func(c *fiber.Ctx) error { return MenuSingleNameHandler(c, db) })
		admin.Post("/generate-table", func(c *fiber.Ctx) error { return GenerateTableHandler(c, db) })
		admin.Post("/generate-multi-table", func(c *fiber.Ctx) error { return GenerateMultiTableHandler(c, db) })
		admin.Post("/generate-module", func(c *fiber.Ctx) error { return CreateModulesHandler(c, db) })
		admin.Post("/execute-flow", func(c *fiber.Ctx) error { return ExecuteFlowHandler(c, db) })
		admin.Post("/down-template", func(c *fiber.Ctx) error { return DownTemplateHandler(c, db) })
		admin.Post("/down-template", func(c *fiber.Ctx) error { return DownTemplateHandler(c, db) })
	}

	media := app.Group("/media")
	media.Use(AuthMiddleware)
	{
		media.Get("/", ListMedia)
		media.Post("/upload", UploadMedia)
		media.Delete("/delete", DeleteMedia)
		media.Post("/rename", RenameMedia)
		media.Post("/folder", CreateFolder)
		media.Get("/preview", PreviewMedia)
	}
}
