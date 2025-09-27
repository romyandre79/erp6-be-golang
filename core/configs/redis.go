package configs

import (
	"context"
	"erp6-be-golang/core/helpers"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *Config) (*redis.Client, error) {
	DB, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		helpers.IsError(err, "Redis DB", true)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       DB,
	})

	// Ping test
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		helpers.IsError(err, "Redis Server", true)
	}
	return rdb, err
}
