package main

import (
	"erp6-be-golang/core/cache"
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"log"
	"strconv"

	_ "erp6-be-golang/docs"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

//go:generate go run generate.go

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
	cache.NewCache()

	log.Print("End Configuration ... ")

	// Init Http
	bodyLimit, _ := strconv.Atoi(configs.ConfigApps.BodyLimit)
	caseSensitive, _ := strconv.ParseBool(configs.ConfigApps.CaseSensitive)
	disableKeepAlive, _ := strconv.ParseBool(configs.ConfigApps.DisableKeepAlive)
	Concurrency, _ := strconv.Atoi(configs.ConfigApps.Concurrency)

	app := fiber.New(fiber.Config{
		AppName:           configs.ConfigApps.AppName,
		BodyLimit:         bodyLimit,
		CaseSensitive:     caseSensitive,
		Concurrency:       Concurrency,
		EnablePrintRoutes: true,
		DisableKeepalive:  disableKeepAlive,
	})

	log.Print("Load Plugin ...")
	plugin.LoadActivePlugins(db, app)
	log.Print("End Load Plugin ...")

	if configs.ConfigApps.SwaggerActive == "true" {
		app.Get("/swagger/*", swagger.New())
	}

	app.Static("/", "./public")

	app.Listen(":" + configs.ConfigApps.AppPort)
}
