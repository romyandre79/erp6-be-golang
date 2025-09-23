package main

import (
	"context"
	"fmt"
	"log"

	"erp6-be-golang/configs"
	"erp6-be-golang/internal/core/logger"
	"erp6-be-golang/internal/core/registry"
	"erp6-be-golang/internal/core/server"
)

func main() {
	// -----------------------
	// Load config
	// -----------------------
	cfg := configs.LoadConfig()

	// -----------------------
	// Init logger
	// -----------------------
	appLogger := logger.NewLogger(cfg.AppLogMode) // "file" / "db" / "remote"
	ctx := context.Background()
	appLogger.Info(ctx, "Starting erp6-be-golang", map[string]interface{}{
		"port": cfg.ServerPort,
		"env":  cfg.AppEnv,
	})

	// -----------------------
	// Init database (GORM)
	// -----------------------
	db, err := configs.InitDatabase(cfg)
	if err != nil {
		appLogger.Error(ctx, "InitDatabase failed", map[string]interface{}{"error": err.Error()})
		log.Fatal(err)
	}
	appLogger.Info(ctx, "Database connected", nil)

	// -----------------------
	// Init redis
	// -----------------------
	rdb, err := configs.InitRedis(cfg)
	if err != nil {
		appLogger.Error(ctx, "InitRedis failed", map[string]interface{}{"error": err.Error()})
		log.Fatal(err)
	}
	appLogger.Info(ctx, "Redis connected", nil)

	// -----------------------
	// Prepare module registry
	// -----------------------
	reg := registry.NewRegistry()
	reg.LoadModules(cfg.Modules)

	if err := reg.InitAll(ctx, db, rdb, appLogger); err != nil {
		log.Fatal("Failed to init modules: ", err)
	}

	// -----------------------
	// Build router (inject registry)
	// -----------------------
	r := server.NewRouter(reg)

	// -----------------------
	// Start server
	// -----------------------
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	appLogger.Info(ctx, "Server starting", map[string]interface{}{"addr": addr})
	if err := r.Run(addr); err != nil {
		appLogger.Error(ctx, "server failed", map[string]interface{}{"error": err.Error()})
		log.Fatal(err)
	}
}
