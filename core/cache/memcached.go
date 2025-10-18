package cache

import (
	"erp6-be-golang/core/configs"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct {
	client *memcache.Client
}

func init() {
	Register("memcached", func() (Cache, error) {
		return NewMemcachedCache()
	})
}

func NewMemcachedCache() (*MemcachedCache, error) {
	addr := configs.ConfigApps.CacheAddr
	if addr == "" {
		addr = "127.0.0.1:11211"
	}
	client := memcache.New(addr)
	return &MemcachedCache{client: client}, nil
}

func (m *MemcachedCache) Set(key string, value interface{}) error {
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(fmt.Sprint(value))})
}

func (m *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}
	return string(item.Value), nil
}

func (m *MemcachedCache) Delete(key string) error {
	return m.client.Delete(key)
}
