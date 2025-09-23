package configs

import (
	"os"
	"strconv"
	"strings"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppLogMode  string // "file" | "db" | "remote" | "std"
	ServerPort  string
	DBDriver    string
	DBSource    string
	RedisAddr   string
	RedisPass   string
	RedisDB     int
	Modules     []string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found, using system env")
	}

	c := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppLogMode: getEnv("APP_LOG_MODE", "std"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBSource:   getEnv("DB_SOURCE", "host=localhost user=postgres password=postgres dbname=erp6 sslmode=disable"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:  getEnv("REDIS_PASS", ""),
		RedisDB:    0,
	}

	if v := getEnv("REDIS_DB", "0"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.RedisDB = n
		}
	}

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
