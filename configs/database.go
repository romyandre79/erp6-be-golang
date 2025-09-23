package configs

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
		case "postgres":
			dialector = postgres.Open(cfg.DBSource)
		case "mysql", "mariadb":
			// expect cfg.DBSource like: user:pass@tcp(localhost:3306)/dbname?parseTime=true
			dialector = mysql.Open(cfg.DBSource)
		case "sqlserver":
			dialector = sqlserver.Open(cfg.DBSource)
		default:
			log.Fatalf("unsupported DB_DRIVER: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}
	return db, nil
}
