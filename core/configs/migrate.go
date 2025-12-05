package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

// RunMigrations runs database migrations if DB_AUTO_MIGRATE is enabled
func RunMigrations(db *gorm.DB) error {
	if strings.ToLower(ConfigApps.DBAutoMigrate) != "true" {
		log.Print("DB Auto-Migration is disabled (DB_AUTO_MIGRATE=false)")
		return nil
	}

	log.Print("Starting database migration...")

	// Determine which SQL file to execute based on database driver
	var sqlFile string
	switch ConfigApps.DBDriver {
	case "mysql", "mariadb":
		sqlFile = "core/db/db-init-maria.sql"
	case "postgres":
		sqlFile = "core/db/db-init-postgres.sql"
	case "sqlserver":
		sqlFile = "core/db/db-init-sqlserver.sql"
	default:
		return fmt.Errorf("unsupported database driver: %s", ConfigApps.DBDriver)
	}

	// Check if SQL file exists
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("SQL migration file not found: %s", sqlFile)
	}

	// Read SQL file
	log.Printf("Reading SQL file: %s", sqlFile)
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Execute SQL statements
	if err := executeSQLFile(db, string(sqlContent)); err != nil {
		return fmt.Errorf("failed to execute SQL migration: %v", err)
	}

	log.Print("✅ Database migration complete")
	return nil
}

// executeSQLFile splits and executes SQL statements from a file
func executeSQLFile(db *gorm.DB, sqlContent string) error {
	// Split SQL content by semicolons to get individual statements
	statements := strings.Split(sqlContent, ";")

	executedCount := 0
	for _, stmt := range statements {
		// Trim whitespace and skip empty statements
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Skip SQL comments
		if strings.HasPrefix(stmt, "--") || strings.HasPrefix(stmt, "/*") {
			continue
		}

		// Execute the statement
		if err := db.Exec(stmt).Error; err != nil {
			// Log warning for non-critical errors (like table already exists)
			if strings.Contains(err.Error(), "already exists") ||
				strings.Contains(err.Error(), "Duplicate") {
				log.Printf("Warning: %v", err)
				continue
			}
			return fmt.Errorf("error executing statement: %v\nStatement: %s", err, stmt[:min(len(stmt), 100)])
		}
		executedCount++
	}

	log.Printf("Executed %d SQL statements successfully", executedCount)
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
