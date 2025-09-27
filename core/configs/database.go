package configs

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
		)
		dialector = postgres.Open(dsn)
	case "mysql", "mariadb":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
		dialector = mysql.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
		dialector = sqlserver.Open(dsn)
	default:
		helpers.IsEmptyLog("Unsupported DB_DRIVER: ", cfg.DBDriver, true)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	helpers.IsError(err, "Check DB Server", true)

	sqlDB, err := db.DB()
	helpers.IsError(err, "Check Open DB", true)

	IdleConv, err := strconv.Atoi(cfg.DBIdleConn)
	if err == nil {
		sqlDB.SetMaxIdleConns(IdleConv)
	}

	MaxConv, err := strconv.Atoi(cfg.DBMaxConn)
	if err == nil {
		sqlDB.SetMaxOpenConns(MaxConv)
	}

	return db, err
}
