package configs

import (
	"erp6-be-golang/core/helpers"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName          string
	AppEnv           string
	AppPort          string
	ReadBufferSize   string
	CaseSensitive    string
	Concurrency      string
	BodyLimit        string
	DBDriver         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	DBIdleConn       string
	DBMaxConn        string
	DisableKeepAlive string
	JwtSecret        string
	LimiterMax       string
	LimiterExpire    string
	LogMode          string
	LogFile          string
	LogRemote        string
	LogDb            string
	CacheType        string
	CacheAddr        string
	CachePass        string
	CacheDB          string
	SwaggerActive    string
	StorageType      string
	WriteBufferSize  string
	Modules          []string
}

var ConfigApps *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found, using system env")
	}

	ConfigApps = &Config{
		AppName:          getEnv("APP_NAME", ""),
		AppEnv:           getEnv("APP_ENV", ""),
		AppPort:          getEnv("APP_PORT", ""),
		BodyLimit:        getEnv("BODY_LIMIT", ""),
		CaseSensitive:    getEnv("CASE_SENSITIVE", ""),
		Concurrency:      getEnv("CONCURRENCY", ""),
		DBDriver:         getEnv("DB_DRIVER", ""),
		DBHost:           getEnv("DB_HOST", ""),
		DBPort:           getEnv("DB_PORT", ""),
		DBUser:           getEnv("DB_USER", ""),
		DBPass:           getEnv("DB_PASS", ""),
		DBName:           getEnv("DB_NAME", ""),
		DBIdleConn:       getEnv("DB_IDLE_CONN", ""),
		DBMaxConn:        getEnv("DB_MAX_CONN", ""),
		DisableKeepAlive: getEnv("DISABLE_KEEP_ALIVE", ""),
		JwtSecret:        getEnv("JWT_SECRET", ""),
		LimiterMax:       getEnv("LIMITER_MAX", ""),
		LimiterExpire:    getEnv("LIMITER_EXPIRE", ""),
		LogMode:          getEnv("LOG_MODE", ""),
		LogFile:          getEnv("LOG_FILE", ""),
		LogRemote:        getEnv("LOG_REMOTE", ""),
		LogDb:            getEnv("LOG_DB", ""),
		CacheType:        getEnv("CACHE_TYPE", ""),
		CacheAddr:        getEnv("CACHE_ADDR", ""),
		CachePass:        getEnv("CACHE_PASS", ""),
		CacheDB:          getEnv("CACHE_DB", ""),
		ReadBufferSize:   getEnv("READ_BUFFER_SIZE", ""),
		SwaggerActive:    getEnv("SWAGGER_ACTIVE", ""),
		StorageType:      getEnv("STORAGE_TYPE", ""),
		WriteBufferSize:  getEnv("WRITE_BUFFER_SIZE", ""),
	}

	// check .env details
	helpers.IsEmptyLog(ConfigApps.AppName, "APP_NAME", true)
	helpers.IsEmptyLog(ConfigApps.AppEnv, "APP_ENV", true)
	helpers.IsEmptyLog(ConfigApps.AppPort, "APP_PORT", true)
	helpers.IsEmptyLog(ConfigApps.BodyLimit, "BODY_LIMIT", true)
	helpers.IsEmptyLog(ConfigApps.CaseSensitive, "CASE_SENSITIVE", true)
	helpers.IsEmptyLog(ConfigApps.Concurrency, "CONCURRENCY", true)
	helpers.IsEmptyLog(ConfigApps.DBDriver, "DB_DRIVER", true)
	helpers.IsEmptyLog(ConfigApps.DBHost, "DB_HOST", true)
	helpers.IsEmptyLog(ConfigApps.DBPort, "DB_PORT", true)
	helpers.IsEmptyLog(ConfigApps.DBUser, "DB_USER", true)
	helpers.IsEmptyLog(ConfigApps.DBPass, "DB_PASS", true)
	helpers.IsEmptyLog(ConfigApps.DBName, "DB_NAME", true)
	helpers.IsEmptyLog(ConfigApps.DBIdleConn, "DB_IDLE_CONN", true)
	helpers.IsEmptyLog(ConfigApps.DBMaxConn, "DB_MAX_CONN", true)
	helpers.IsEmptyLog(ConfigApps.DisableKeepAlive, "DISABLE_KEEP_ALIVE", true)
	helpers.IsEmptyLog(ConfigApps.LogMode, "LOG_MODE", true)
	helpers.IsEmptyLog(ConfigApps.LogFile, "LOG_FILE", false)
	helpers.IsEmptyLog(ConfigApps.LogRemote, "LOG_REMOTE", false)
	helpers.IsEmptyLog(ConfigApps.LogDb, "LOG_DB", true)
	helpers.IsEmptyLog(ConfigApps.CacheType, "CACHE_TYPE", false)
	helpers.IsEmptyLog(ConfigApps.CacheAddr, "REDIS_ADDR", false)
	helpers.IsEmptyLog(ConfigApps.CachePass, "REDIS_PASS", false)
	helpers.IsEmptyLog(ConfigApps.CacheDB, "REDIS_DB", true)
	helpers.IsEmptyLog(ConfigApps.LimiterMax, "LIMITER_MAX", true)
	helpers.IsEmptyLog(ConfigApps.LimiterExpire, "LIMITER_EXPIRE", true)
	helpers.IsEmptyLog(ConfigApps.ReadBufferSize, "READ_BUFFER_SIZE", true)
	helpers.IsEmptyLog(ConfigApps.SwaggerActive, "SWAGGER_ACTIVE", true)
	helpers.IsEmptyLog(ConfigApps.StorageType, "STORAGE_TYPE", false)
	helpers.IsEmptyLog(ConfigApps.WriteBufferSize, "WRITE_BUFFER_SIZE", true)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
