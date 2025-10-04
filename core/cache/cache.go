package cache

import (
	"erp6-be-golang/core/configs"
	"fmt"
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

var registry = map[string]func() (Cache, error){}

// Register backend ke registry
func Register(name string, factory func() (Cache, error)) {
	registry[name] = factory
}

// Factory pilih backend berdasarkan ENV
func NewCache() (Cache, error) {
	backend := configs.ConfigApps.CacheType // contoh: "redis", "memory"
	if backend != "" {
		if factory, ok := registry[backend]; ok {
			return factory()
		}
		return nil, fmt.Errorf("unknown cache backend: %s", backend)
	} else {
		return nil, nil
	}
}
