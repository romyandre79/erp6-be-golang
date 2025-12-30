package main

import (
	"erp6-be-golang/core/cache"
	"erp6-be-golang/core/configs"
	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/core/scheduler"
	"time"

	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/models"
	_ "erp6-be-golang/plugins/admin"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	// Initialize WA DB for AI
	generator.SetWADatabase(db)
	
	// Start WA Client immediately
	go func() {
		if err := generator.InitWhatmeow(); err != nil {
			log.Printf("Failed to init WhatsApp: %v", err)
		}
	}()

	// Initialize Telegram DB
	generator.SetTelegramDatabase(db)
	
	// Start Telegram Bot if token is configured
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	log.Printf("Checking Telegram configuration... Token length: %d", len(telegramToken))
	
	if telegramToken != "" {
		log.Println("Starting Telegram Bot...")
		go func() {
			if err := generator.InitTelegram(telegramToken); err != nil {
				log.Printf("Failed to init Telegram: %v", err)
			} else {
				log.Println("Telegram Bot initialized successfully!")
			}
		}()
	} else {
		log.Println("Telegram bot token NOT configured in .env (TELEGRAM_BOT_TOKEN is empty)")
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
	prefork, _ := strconv.ParseBool(configs.ConfigApps.Prefork)
	EnablePrintRoutes, _ := strconv.ParseBool(configs.ConfigApps.EnablePrintRoutes)
	disableKeepAlive, _ := strconv.ParseBool(configs.ConfigApps.DisableKeepAlive)
	Concurrency, _ := strconv.Atoi(configs.ConfigApps.Concurrency)
	ReadBufferSize, _ := strconv.Atoi(configs.ConfigApps.ReadBufferSize)
	WriteBufferSize, _ := strconv.Atoi(configs.ConfigApps.WriteBufferSize)
	LimiterMax, _ := strconv.Atoi(configs.ConfigApps.LimiterMax)
	LimiterExpire, _ := strconv.Atoi(configs.ConfigApps.LimiterExpire)

	app := fiber.New(fiber.Config{
		AppName:           configs.ConfigApps.AppName,
		BodyLimit:         bodyLimit,
		CaseSensitive:     caseSensitive,
		Concurrency:       Concurrency,
		EnablePrintRoutes: EnablePrintRoutes,
		DisableKeepalive:  disableKeepAlive,
		ReadBufferSize:    ReadBufferSize,
		WriteBufferSize:   WriteBufferSize,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		Prefork:           prefork,
	})

	app.Use(limiter.New(limiter.Config{
		Max:        LimiterMax,
		Expiration: time.Duration(LimiterExpire),
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit per IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Enable Gzip/Brotli Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance speed vs size
	}))

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

	app.Use(func(c *fiber.Ctx) error {
		// Serve static files from public directory
		path := c.Path()
		if !strings.HasPrefix(path, "/api/") {
			// Try to serve static file
			filePath := filepath.Join("./public", path)
			if _, err := os.Stat(filePath); err == nil {
				return c.SendFile(filePath)
			}
		}
		return c.Next()
	})

	app.Get("*", func(c *fiber.Ctx) error {
		// Only serve index.html for non-API routes and non-static files
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
