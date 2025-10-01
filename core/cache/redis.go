package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func init() {
	Register("redis", func() (Cache, error) {
		return NewRedisCache()
	})
}

func NewRedisCache() (*RedisCache, error) {
	db := configs.ConfigApps.CacheDB
	client := redis.NewClient(&redis.Options{
		Addr:     configs.ConfigApps.CacheAddr
		Password: configs.ConfigApps.CachePass
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Set(key string, value interface{}) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
