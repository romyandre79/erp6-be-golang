package configs

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// EnsureDatabaseExists checks if the database exists and creates it if it doesn't
func EnsureDatabaseExists() error {
	var db *gorm.DB
	var err error
	var checkQuery string
	var createQuery string

	switch ConfigApps.DBDriver {
	case "postgres":
		// Connect to postgres database (default database)
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			ConfigApps.DBHost, ConfigApps.DBPort, ConfigApps.DBUser, ConfigApps.DBPass,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL server: %v", err)
		}
		checkQuery = fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", ConfigApps.DBName)
		createQuery = fmt.Sprintf("CREATE DATABASE %s", ConfigApps.DBName)

	case "mysql", "mariadb":
		// Connect without specifying database
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/?parseTime=true",
			ConfigApps.DBUser, ConfigApps.DBPass, ConfigApps.DBHost, ConfigApps.DBPort,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to MySQL/MariaDB server: %v", err)
		}
		checkQuery = fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", ConfigApps.DBName)
		createQuery = fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", ConfigApps.DBName)

	case "sqlserver":
		// Connect to master database
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=master",
			ConfigApps.DBUser, ConfigApps.DBPass, ConfigApps.DBHost, ConfigApps.DBPort,
		)
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to SQL Server: %v", err)
		}
		checkQuery = fmt.Sprintf("SELECT name FROM sys.databases WHERE name = '%s'", ConfigApps.DBName)
		createQuery = fmt.Sprintf("CREATE DATABASE [%s]", ConfigApps.DBName)

	default:
		return fmt.Errorf("unsupported DB_DRIVER: %s", ConfigApps.DBDriver)
	}

	// Check if database exists
	var result string
	err = db.Raw(checkQuery).Scan(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check database existence: %v", err)
	}

	// Create database if it doesn't exist
	if result == "" {
		log.Printf("Database '%s' does not exist. Creating...", ConfigApps.DBName)
		err = db.Exec(createQuery).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("✅ Database '%s' created successfully", ConfigApps.DBName)
	} else {
		log.Printf("Database '%s' already exists", ConfigApps.DBName)
	}

	// Close the connection
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	return nil
}

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
