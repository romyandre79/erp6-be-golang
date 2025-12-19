// admin/auth.go
package admin

import (
	"encoding/json"
	dbgenerator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/models"

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
	ws.GlobalHub = ws.NewHub()
	go ws.GlobalHub.Run()

	// Implement Hub Callbacks
	ws.GlobalHub.OnUserOnline = func(userID int) {
		db.Model(&models.Useraccess{}).Where("useraccessid = ?", userID).Update("isonline", 1)
		msg, _ := json.Marshal(map[string]interface{}{
			"type":     "status_update",
			"user_id":  userID,
			"isonline": 1,
		})
		ws.GlobalHub.Broadcast <- msg
	}

	ws.GlobalHub.OnUserOffline = func(userID int) {
		db.Model(&models.Useraccess{}).Where("useraccessid = ?", userID).Update("isonline", 0)
		msg, _ := json.Marshal(map[string]interface{}{
			"type":     "status_update",
			"user_id":  userID,
			"isonline": 0,
		})
		ws.GlobalHub.Broadcast <- msg
	}

	// Set GlobalDB for WebSocket handlers
	GlobalDB = db

	// Public routes
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/login", func(c *fiber.Ctx) error { return LoginHandler(c, db) })
	auth.Post("/load-theme", func(c *fiber.Ctx) error { return LoadThemeHandler(c, db) })

	// Public CMS Routes
	public := app.Group("/api/public")
	public.Get("/page/*", func(c *fiber.Ctx) error { return PublicPageHandler(c, db) })

	// Protected routes
	auth.Use(AuthMiddleware)
	{
		//auth
		auth.Post("/logout", func(c *fiber.Ctx) error { return LogoutHandler(c, db) })
		auth.Get("/me", func(c *fiber.Ctx) error { return MeHander(c, db) })
	}

	admin := app.Group("/api/admin")
	admin.Use(AuthMiddleware)
	{
		admin.Get("/getmenu", func(c *fiber.Ctx) error { return MenuSingleNameHandler(c, db) })
		admin.Post("/generate-table", func(c *fiber.Ctx) error { return GenerateTableHandler(c, db) })
		admin.Post("/generate-multi-table", func(c *fiber.Ctx) error { return GenerateMultiTableHandler(c, db) })
		admin.Post("/generate-module", func(c *fiber.Ctx) error { return CreateModulesHandler(c, db) })
		admin.Post("/execute-flow", func(c *fiber.Ctx) error { return ExecuteFlowHandler(c, db) })
		admin.Post("/down-template", func(c *fiber.Ctx) error { return DownTemplateHandler(c, db) })
		admin.Post("/execute-table-operation", func(c *fiber.Ctx) error { return ExecuteTableOperationHandler(c, db) })
		admin.Post("/ai/command", func(c *fiber.Ctx) error { return AiCommandHandler(c, db) })
		admin.Post("/plugins/upload", func(c *fiber.Ctx) error {
			return dbgenerator.HandlePluginUpload(c, db)
		})

		// Module Management Routes
		admin.Post("/module/upload", func(c *fiber.Ctx) error {
			return UploadModulePackageHandler(c, db)
		})
		admin.Post("/module/uninstall/:moduleid", func(c *fiber.Ctx) error {
			return UninstallModuleHandler(c, db)
		})
		admin.Get("/module/details/:moduleid", func(c *fiber.Ctx) error {
			return GetModuleDetailsHandler(c, db)
		})
		admin.Get("/module/dependencies/:moduleid", func(c *fiber.Ctx) error {
			return GetModuleDependenciesHandler(c, db)
		})
		admin.Get("/module/export/:moduleid", func(c *fiber.Ctx) error {
			return ExportModuleHandler(c, db)
		})

		// Notification Routes
		admin.Get("/notifications/unread", func(c *fiber.Ctx) error { return GetUnreadNotifications(c, db) })
		admin.Post("/notifications/:id/read", func(c *fiber.Ctx) error { return MarkAsRead(c, db) })
		admin.Post("/notifications/send", func(c *fiber.Ctx) error { return SendNotification(c, db) })

		admin.Get("/chat/history", func(c *fiber.Ctx) error { return GetChatHistoryHandler(c, db) })

		// DB Backup/Restore
		admin.Post("/db/backup", func(c *fiber.Ctx) error { return BackupHandler(c, db) })
		admin.Post("/db/restore", func(c *fiber.Ctx) error { return RestoreHandler(c, db) })

		// DB Reverse Engineering
		admin.Post("/db/reverse-engineer", func(c *fiber.Ctx) error { return ReverseEngineerHandler(c, db) })

		// DB Resources (Relations & Areas)
		admin.Post("/db/relations/save", func(c *fiber.Ctx) error { return SaveRelationsHandler(c, db) })
		admin.Get("/db/relations", func(c *fiber.Ctx) error { return GetRelationsHandler(c, db) })
		admin.Delete("/db/relations/:id", func(c *fiber.Ctx) error { return DeleteRelationHandler(c, db) })

		admin.Post("/db/areas/save", func(c *fiber.Ctx) error { return SaveAreasHandler(c, db) })
		admin.Get("/db/areas", func(c *fiber.Ctx) error { return GetAreasHandler(c, db) })
		admin.Delete("/db/areas/:id", func(c *fiber.Ctx) error { return DeleteAreaHandler(c, db) })

		// Report Designer Routes
		admin.Post("/report-templates", func(c *fiber.Ctx) error { return CreateReport(c, db) })
		admin.Get("/report-templates", func(c *fiber.Ctx) error { return ListReports(c, db) })
		admin.Get("/report-templates/:id", func(c *fiber.Ctx) error { return GetReport(c, db) })
		admin.Put("/report-templates/:id", func(c *fiber.Ctx) error { return UpdateReport(c, db) })
		admin.Delete("/report-templates/:id", func(c *fiber.Ctx) error { return DeleteReport(c, db) })
		admin.Post("/report-templates/:id/preview", func(c *fiber.Ctx) error { return PreviewReport(c, db) })
		admin.Post("/report-templates/:id/execute", func(c *fiber.Ctx) error { return ExecuteReport(c, db) })

		// JRXML Import/Export Routes
		admin.Post("/report-templates/import-jrxml", func(c *fiber.Ctx) error { return ImportJRXML(c, db) })
		admin.Get("/report-templates/:id/export-jrxml", func(c *fiber.Ctx) error { return ExportJRXML(c, db) })

		// Scheduler Management
		admin.Post("/scheduler/reload", func(c *fiber.Ctx) error { return ReloadSchedulerHandler(c, db) })

		// Record Locking
		admin.Post("/lock-record", func(c *fiber.Ctx) error { return LockRecordHandler(c, db) })
		admin.Post("/unlock-record", func(c *fiber.Ctx) error { return UnlockRecordHandler(c, db) })
	}

	app.Get("/api/ws/notifications",
		AuthMiddleware, // ✅ Auth first
		func(c *fiber.Ctx) error { // ✅ Then check WS upgrade
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		},
		websocket.New(WebSocketHandler),
	)

	media := app.Group("/api/media")
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
