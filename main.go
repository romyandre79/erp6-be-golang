package main

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/logger"
	"erp6-be-golang/core/plugin"
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

//go:generate go run generate.go

func main() {

	// Load Config from .env
	log.Print("Check Configuration ... ")
	configApps := configs.LoadConfig()

	// Load Logger from .env
	_, err := logger.InitLogger(configApps)

	if err != nil {
		helpers.IsError(err, "Init Logger", true)
	}

	// Load DB from .env
	db, err := configs.InitDatabase(configApps)

	if err != nil {
		helpers.IsError(err, "Check DB Server", true)
	}

	// Load Redis from .env
	_, err = configs.InitRedis(configApps)
	if err != nil {
		helpers.IsError(err, "Redis Server", true)
	}

	log.Print("End Configuration ... ")

	// Init Http
	app := fiber.New(fiber.Config{})

	log.Print("Load Plugin ...")
	plugin.LoadActivePlugins(db, app)
	log.Print("End Load Plugin ...")

	app.Use(swagger.New())

	app.Listen(":" + configApps.AppPort)
}
