package main

import (
	"erp6-be-golang/core/cache"
	"erp6-be-golang/core/configs"
	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/core/scheduler"

	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/i18n"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"flag"
	"log"
	"strconv"
	"strings"
	"erp6-be-golang/models"
	_ "erp6-be-golang/plugins/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"fmt"
)

//go:generate go run generate.go

// @title ERP6 API
// @version 1.0
// @description API ERP6 dengan Fiber Swagger
// @host localhost:8888
// @BasePath /
func main() {
	// CLI Flags
	backupPtr := flag.Bool("backup", false, "Backup database to file (default: backup.sql or backup.db)")
	restorePtr := flag.String("restore", "", "Restore database from file")
	flag.Parse()

	i18n.Init()

	// Load Config from .env
	log.Print("Check Configuration ... ")
	configs.LoadConfig()

	// Handle CLI Operations
	if *backupPtr {
		log.Println("Starting Database Backup...")
		outFile := "backup.sql"
		if configs.ConfigApps.DBDriver == "sqlite" || configs.ConfigApps.DBDriver == "sqlite3" {
			outFile = "backup.db"
		}
		err := generator.BackupDatabase(
			configs.ConfigApps.DBDriver,
			configs.ConfigApps.DBHost,
			configs.ConfigApps.DBPort,
			configs.ConfigApps.DBUser,
			configs.ConfigApps.DBPass,
			configs.ConfigApps.DBName,
			outFile,
		)
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		log.Printf("Backup successful! File saved to: %s", outFile)
		return
	}

	if *restorePtr != "" {
		log.Printf("Starting Database Restore from %s...", *restorePtr)
		err := generator.RestoreDatabase(
			configs.ConfigApps.DBDriver,
			configs.ConfigApps.DBHost,
			configs.ConfigApps.DBPort,
			configs.ConfigApps.DBUser,
			configs.ConfigApps.DBPass,
			configs.ConfigApps.DBName,
			*restorePtr,
		)
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
		log.Println("Restore successful!")
		return
	}

	// Load Logger from .env
	_, err := logger.InitLogger()

	if err != nil {
		helpers.IsError(err, "Init Logger", true)
	}

	// Load DB from .env
	// Ensure database exists (create if not)
	if err := configs.EnsureDatabaseExists(); err != nil {
		helpers.IsError(err, "Ensure Database Exists", true)
	}

	db, err := configs.InitDatabase()

	if err != nil {
		helpers.IsError(err, "Check DB Server", true)
	}

	// Run database migrations if enabled
	// Auto Migrate
	db.AutoMigrate(&models.Chat{})
	err = configs.RunMigrations(db)
	if err != nil {
		helpers.IsError(err, "Database Migration", true)
	}

	// Load Cache from .env
	_, err = cache.NewCache()
	if err != nil {
		helpers.IsError(err, "Cache Server", true)
	}

	log.Print("End Configuration ... ")

	// Init Http
	bodyLimit, _ := strconv.Atoi(configs.ConfigApps.BodyLimit)
	caseSensitive, _ := strconv.ParseBool(configs.ConfigApps.CaseSensitive)
	disableKeepAlive, _ := strconv.ParseBool(configs.ConfigApps.DisableKeepAlive)
	Concurrency, _ := strconv.Atoi(configs.ConfigApps.Concurrency)
	ReadBufferSize, _ := strconv.Atoi(configs.ConfigApps.ReadBufferSize)
	WriteBufferSize, _ := strconv.Atoi(configs.ConfigApps.WriteBufferSize)

	app := fiber.New(fiber.Config{
		AppName:           configs.ConfigApps.AppName,
		BodyLimit:         bodyLimit,
		CaseSensitive:     caseSensitive,
		Concurrency:       Concurrency,
		EnablePrintRoutes: true,
		DisableKeepalive:  disableKeepAlive,
		ReadBufferSize:    ReadBufferSize,
		WriteBufferSize:   WriteBufferSize,
	})

	allowedOrigins := strings.Split(configs.ConfigApps.AllowOrigin, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	log.Print("Load Plugin ...")
	plugin.LoadActivePlugins(db, app)
	log.Print("End Load Plugin ...")

	log.Print("Load Workflow Components ...")
	generator.LoadPlugins(db)
	log.Print("End Load Workflow Components ...")

	app.All("/api/webhook/:source", func(c *fiber.Ctx) error {
		return plugin.WebhookHandler(c, db)
	})

	app.Static("/", "./public")

	// Catch-all route for SPA (must be last)
	app.Post("/api/public/request-demo", func(c *fiber.Ctx) error { // Request Demo Handler
		type DemoRequest struct {
			Name              string   `json:"name"`
			Email             string   `json:"email"`
			Phone             string   `json:"phone"`
			SelectedApps      []string `json:"selectedApps"`
			SelectedWorkflows []string `json:"selectedWorkflows"`
			UsersCount        int      `json:"usersCount"`
			StorageSize       int      `json:"storageSize"`
			DeploymentMode    string   `json:"deploymentMode"`
			ConsultationHours int      `json:"consultationHours"`
			TransportDays     int      `json:"transportDays"`
			MonthlyRecurring  float64  `json:"monthlyRecurring"`
			OneTimeFee        float64  `json:"oneTimeFee"`
		}

		var req DemoRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		// Send email
		to := []string{"romyandre79@gmail.com"}
		subject := "New Demo Request from " + req.Name

		// Construct formatted lists
		appsList := "None"
		if len(req.SelectedApps) > 0 {
			appsList = fmt.Sprintf("%v", req.SelectedApps)
		}
		workflowsList := "None"
		if len(req.SelectedWorkflows) > 0 {
			workflowsList = fmt.Sprintf("%v", req.SelectedWorkflows)
		}

		// EMAIL TEMPLATE - You can adjust the email body here
		emailTemplate := `
Hello Admin,

You have received a new demo request for the ERP6 System.

User Details:
-------------
Name:  %s
Email: %s
Phone: %s

Configuration:
--------------
Deployment: %s
Users: %d
Storage: %d GB
Apps: %s
Workflows: %s

Services:
---------
Consultation: %d Hours
Transport: %d Days

Est. Pricing:
-------------
Monthly Recurring: Rp %.0f
One-Time Fee: Rp %.0f

Please contact them shortly.

Best Regards,
ERP6 Auto-Bot
`
		body := fmt.Sprintf(emailTemplate,
			req.Name, req.Email, req.Phone,
			req.DeploymentMode, req.UsersCount, req.StorageSize, appsList, workflowsList,
			req.ConsultationHours, req.TransportDays,
			req.MonthlyRecurring, req.OneTimeFee,
		)

		if err := helpers.SendEmail(to, subject, body); err != nil {
			// Log error but maybe don't fail the request to the user if it's just email failure?
			// Or arguably we should tell them. Let's return error for now so they know.
			// Converting err to string for simple logging
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to send email: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Demo request sent successfully",
		})
	})

	app.Get("*", func(c *fiber.Ctx) error {
		// Only serve index.html for non-API routes
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendFile("./public/index.html")
	})

	// Init Scheduler (after App creation)
	scheduler.Init(app)
	scheduler.LoadJobs(db)

	if configs.ConfigApps.IsHttps == "true" {
		log.Print("Server listening on HTTPS ...")
		if err := app.ListenTLS(configs.ConfigApps.AppHost+":"+configs.ConfigApps.AppPort, configs.ConfigApps.SslCert, configs.ConfigApps.SslKey); err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		log.Print("Server listening on HTTP ...")
		app.Listen(configs.ConfigApps.AppHost + ":" + configs.ConfigApps.AppPort)
	}
}
