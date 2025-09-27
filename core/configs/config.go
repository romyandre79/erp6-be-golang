package configs

import (
	"erp6-be-golang/core/helpers"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppPort    string
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBIdleConn string
	DBMaxConn  string
	LogMode    string
	LogFile    string
	LogRemote  string
	LogDb      string
	RedisAddr  string
	RedisPass  string
	RedisDB    string
	Modules    []string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found, using system env")
	}

	c := &Config{
		AppEnv:     getEnv("APP_ENV", ""),
		AppPort:    getEnv("APP_PORT", ""),
		DBDriver:   getEnv("DB_DRIVER", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPass:     getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBIdleConn: getEnv("DB_IDLE_CONN", ""),
		DBMaxConn:  getEnv("DB_MAX_CONN", "1"),
		LogMode:    getEnv("LOG_MODE", ""),
		LogFile:    getEnv("LOG_FILE", ""),
		LogRemote:  getEnv("LOG_REMOTE", ""),
		LogDb:      getEnv("LOG_DB", ""),
		RedisAddr:  getEnv("REDIS_ADDR", ""),
		RedisPass:  getEnv("REDIS_PASS", ""),
		RedisDB:    getEnv("REDIS_DB", ""),
	}

	// check .env details
	helpers.IsEmptyLog(c.AppEnv, "APP_ENV", true)
	helpers.IsEmptyLog(c.AppPort, "APP_PORT", true)
	helpers.IsEmptyLog(c.DBDriver, "DB_DRIVER", true)
	helpers.IsEmptyLog(c.DBHost, "DB_HOST", true)
	helpers.IsEmptyLog(c.DBPort, "DB_PORT", true)
	helpers.IsEmptyLog(c.DBUser, "DB_USER", true)
	helpers.IsEmptyLog(c.DBPass, "DB_PASS", true)
	helpers.IsEmptyLog(c.DBName, "DB_NAME", true)
	helpers.IsEmptyLog(c.DBIdleConn, "DB_IDLE_CONN", true)
	helpers.IsEmptyLog(c.DBMaxConn, "DB_MAX_CONN", true)
	helpers.IsEmptyLog(c.LogMode, "LOG_MODE", true)
	helpers.IsEmptyLog(c.LogFile, "LOG_FILE", false)
	helpers.IsEmptyLog(c.LogRemote, "LOG_REMOTE", false)
	helpers.IsEmptyLog(c.LogDb, "LOG_DB", true)
	helpers.IsEmptyLog(c.RedisAddr, "REDIS_ADDR", true)
	helpers.IsEmptyLog(c.RedisPass, "REDIS_PASS", false)
	helpers.IsEmptyLog(c.RedisDB, "REDIS_DB", true)

	// optional MODULES env (comma separated)
	if m := getEnv("MODULES", ""); m != "" {
		for _, s := range strings.Split(m, ",") {
			c.Modules = append(c.Modules, strings.TrimSpace(s))
		}
	}
	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
