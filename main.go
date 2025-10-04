package main

import (
	"erp6-be-golang/core/cache"
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/email"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"erp6-be-golang/core/storage"
	"log"
	"strconv"

	_ "erp6-be-golang/docs"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

//go:generate go run generate.go

// @title ERP6 API
// @version 1.0
// @description API ERP6 dengan Fiber Swagger
// @host localhost:8888
// @BasePath /
func main() {

	// Load Config from .env
	log.Print("Check Configuration ... ")
	configs.LoadConfig()

	// Load Logger from .env
	_, err := logger.InitLogger()

	if err != nil {
		helpers.IsError(err, "Init Logger", true)
	}

	// Load DB from .env
	db, err := configs.InitDatabase()

	if err != nil {
		helpers.IsError(err, "Check DB Server", true)
	}

	// Load Cache from .env
	_, err = cache.NewCache()
	if err != nil {
		helpers.IsError(err, "Cache Server", true)
	}

	// Load Storage form .env
	_, error := storage.Get()
	if error {
		helpers.IsError(err, "Storage Server", true)
	}

	_, err = email.NewEmailSender()
	if err != nil {
		helpers.IsError(err, "Email Server", true)
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

	log.Print("Load Plugin ...")
	plugin.LoadActivePlugins(db, app)
	log.Print("End Load Plugin ...")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		BasePath: "/", // default
	}))

	app.Static("/", "./public")

	app.Listen(":" + configs.ConfigApps.AppPort)
}
