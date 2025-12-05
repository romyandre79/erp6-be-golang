// admin/auth.go
package admin

import (
	dbgenerator "erp6-be-golang/core/generator/db"

	"erp6-be-golang/core/ws"

	"github.com/gofiber/contrib/websocket"
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
	// Initialize Hub
	Hub = ws.NewHub()
	go Hub.Run()

	auth := app.Group("/auth")

	// Public routes
	auth.Post("/login", func(c *fiber.Ctx) error { return LoginHandler(c, db) })
	auth.Post("/load-theme", func(c *fiber.Ctx) error { return LoadThemeHandler(c, db) })

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
		admin.Post("/execute-table-operation", func(c *fiber.Ctx) error { return ExecuteTableOperationHandler(c, db) })
		admin.Post("/plugins/upload", func(c *fiber.Ctx) error {
			return dbgenerator.HandlePluginUpload(c, db)
		})

		// Notification Routes
		admin.Get("/notifications/unread", func(c *fiber.Ctx) error { return GetUnreadNotifications(c, db) })
		admin.Post("/notifications/:id/read", func(c *fiber.Ctx) error { return MarkAsRead(c, db) })
		admin.Post("/notifications/send", func(c *fiber.Ctx) error { return SendNotification(c, db) })
	}

	// WebSocket Route - Specific middleware for WS upgrade might be needed if Headers are missing
	// Just reusing AuthMiddleware for now assuming token is passed in Query param or Header
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Apply Auth Middleware to /ws?
	// If AuthMiddleware reads Authorization header, Browser JS WebSocket cannot set it easily.
	// We might need a query param adapter.
	app.Get("/ws/notifications", AuthMiddleware, websocket.New(WebSocketHandler))

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
