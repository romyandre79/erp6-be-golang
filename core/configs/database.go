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

func InitDatabase() (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch ConfigApps.DBDriver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			ConfigApps.DBHost, ConfigApps.DBPort, ConfigApps.DBUser, ConfigApps.DBPass, ConfigApps.DBName,
		)
		dialector = postgres.Open(dsn)
	case "mysql", "mariadb":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			ConfigApps.DBUser, ConfigApps.DBPass, ConfigApps.DBHost, ConfigApps.DBPort, ConfigApps.DBName,
		)
		dialector = mysql.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			ConfigApps.DBUser, ConfigApps.DBPass, ConfigApps.DBHost, ConfigApps.DBPort, ConfigApps.DBName,
		)
		dialector = sqlserver.Open(dsn)
	default:
		helpers.IsEmptyLog("Unsupported DB_DRIVER: ", ConfigApps.DBDriver, true)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	helpers.IsError(err, "Check DB Server", true)

	sqlDB, err := db.DB()
	helpers.IsError(err, "Check Open DB", true)

	IdleConv, err := strconv.Atoi(ConfigApps.DBIdleConn)
	if err == nil {
		sqlDB.SetMaxIdleConns(IdleConv)
	}

	MaxConv, err := strconv.Atoi(ConfigApps.DBMaxConn)
	if err == nil {
		sqlDB.SetMaxOpenConns(MaxConv)
	}

	return db, err
}
