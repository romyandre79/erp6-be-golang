package admin

import (
	"erp6-be-golang/core/configs"
	generator "erp6-be-golang/core/generator"
	"erp6-be-golang/core/helpers"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BackupHandler(c *fiber.Ctx, dbConn *gorm.DB) error {
	// Generate a unique filename
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)
	if configs.ConfigApps.DBDriver == "sqlite" || configs.ConfigApps.DBDriver == "sqlite3" {
		filename = fmt.Sprintf("backup_%s.db", timestamp)
	}

	backupDir := "./tmp/backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.MkdirAll(backupDir, 0755)
	}

	outputPath := filepath.Join(backupDir, filename)

	err := generator.BackupDatabase(
		configs.ConfigApps.DBDriver,
		configs.ConfigApps.DBHost,
		configs.ConfigApps.DBPort,
		configs.ConfigApps.DBUser,
		configs.ConfigApps.DBPass,
		configs.ConfigApps.DBName,
		outputPath,
	)

	if err != nil {
		return helpers.FailResponse(c, 500, "Backup failed: "+err.Error(), "")
	}

	// Return file download
	return c.Download(outputPath)
}

func RestoreHandler(c *fiber.Ctx, dbConn *gorm.DB) error {
	file, err := c.FormFile("file")
	if err != nil {
		return helpers.FailResponse(c, 400, "File is required", "")
	}

	restoreDir := "./tmp/restores"
	if _, err := os.Stat(restoreDir); os.IsNotExist(err) {
		os.MkdirAll(restoreDir, 0755)
	}

	inputPath := filepath.Join(restoreDir, file.Filename)
	if err := c.SaveFile(file, inputPath); err != nil {
		return helpers.FailResponse(c, 500, "Failed to save file", "")
	}

	err = generator.RestoreDatabase(
		configs.ConfigApps.DBDriver,
		configs.ConfigApps.DBHost,
		configs.ConfigApps.DBPort,
		configs.ConfigApps.DBUser,
		configs.ConfigApps.DBPass,
		configs.ConfigApps.DBName,
		inputPath,
	)

	if err != nil {
		return helpers.FailResponse(c, 500, "Restore failed: "+err.Error(), "")
	}

	return helpers.SuccessResponse(c, "Restore successful", "")
}
