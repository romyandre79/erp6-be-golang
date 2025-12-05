package main

import (
	"erp6-be-golang/core/cache"
	"erp6-be-golang/core/configs"

	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/i18n"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"log"
	"strconv"

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
	i18n.Init()

	// Load Config from .env
	log.Print("Check Configuration ... ")
	configs.LoadConfig()

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

	app.Static("/", "./public")

	app.Listen(":" + configs.ConfigApps.AppPort)
}
