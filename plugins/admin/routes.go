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
		plugin.RegisterModelRoutes(admin, db, models.Language{}, "language")
		plugin.RegisterModelRoutes(admin, db, models.Country{}, "country")
		plugin.RegisterModelRoutes(admin, db, models.Province{}, "province")
		plugin.RegisterModelRoutes(admin, db, models.City{}, "city")
		plugin.RegisterModelRoutes(admin, db, models.Modules{}, "modules")
		plugin.RegisterModelRoutes(admin, db, models.Currency{}, "currency")
		plugin.RegisterModelRoutes(admin, db, models.Theme{}, "theme")
		plugin.RegisterModelRoutes(admin, db, models.Parameter{}, "parameter")
		plugin.RegisterModelRoutes(admin, db, models.Widget{}, "widget")
		plugin.RegisterModelRoutes(admin, db, models.Menuaccess{}, "menuaccess")
		plugin.RegisterModelRoutes(admin, db, models.Groupmenu{}, "groupmenu")
		plugin.RegisterModelRoutes(admin, db, models.Menuauth{}, "menuauth")
		plugin.RegisterModelRoutes(admin, db, models.Groupmenuauth{}, "groupmenuauth")
		plugin.RegisterModelRoutes(admin, db, models.Useraccess{}, "useraccess")
		plugin.RegisterModelRoutes(admin, db, models.Usergroup{}, "usergroup")
		plugin.RegisterModelRoutes(admin, db, models.Snro{}, "snro")
		plugin.RegisterModelRoutes(admin, db, models.Snrodet{}, "snrodet")
		plugin.RegisterModelRoutes(admin, db, models.Workflow{}, "workflow")
		plugin.RegisterModelRoutes(admin, db, models.Wfstatus{}, "wfstatus")
		plugin.RegisterModelRoutes(admin, db, models.Wfgroup{}, "wfgroup")
		plugin.RegisterModelRoutes(admin, db, models.Translog{}, "translog")
		plugin.RegisterModelRoutes(admin, db, models.Apps{}, "apps")
		plugin.RegisterModelRoutes(admin, db, models.Appscompany{}, "appscompany")
		plugin.RegisterModelRoutes(admin, db, models.Appsmenu{}, "appsmenu")

		admin.Get("/getwidget", func(c *fiber.Ctx) error { return DashboardListHandler(c, db) })

	}
}
