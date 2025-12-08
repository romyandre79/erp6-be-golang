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

	_ "erp6-be-golang/plugins/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     configs.ConfigApps.AllowOrigin, // URL Nuxt
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

	app.All("/api/webhook/:source", plugin.WebhookHandler)

	app.Static("/", "./public")

	// Catch-all route for SPA (must be last)
	app.Get("*", func(c *fiber.Ctx) error {
		// Only serve index.html for non-API routes
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendFile("./public/index.html")
	})

	// Init Scheduler (after App creation)
	scheduler.Init(app)
	// Example handler
	scheduler.RegisterHandler("test_job", func() {
		log.Println("Hello from test_job!")
	})
	scheduler.LoadJobs(db)

	app.Listen(":" + configs.ConfigApps.AppPort)
}
