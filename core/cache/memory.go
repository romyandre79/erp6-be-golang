package cache

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type MemoryCache struct {
	c *gocache.Cache
}

func init() {
	Register("memory", func() (Cache, error) {
		return NewMemoryCache(), nil
	})
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		c: gocache.New(5*time.Minute, 10*time.Minute),
	}
}

func (m *MemoryCache) Set(key string, value interface{}) error {
	m.c.Set(key, value, gocache.DefaultExpiration)
	return nil
}

func (m *MemoryCache) Get(key string) (interface{}, error) {
	val, found := m.c.Get(key)
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	return val, nil
}

func (m *MemoryCache) Delete(key string) error {
	m.c.Delete(key)
	return nil
}
