package cache

import (
	"encoding/json"
	"erp6-be-golang/core/configs"
	"os"
	"path/filepath"
)

type FileCache struct {
	path string
}

func init() {
	Register("file", func() (Cache, error) {
		return NewFileCache()
	})
}

func NewFileCache() (*FileCache, error) {
	path := configs.ConfigApps.CacheAddr
	if path == "" {
		path = "./filecache"
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &FileCache{path: path}, nil
}

func (f *FileCache) filePath(key string) string {
	return filepath.Join(f.path, key+".json")
}

func (f *FileCache) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(f.filePath(key), data, 0644)
}

func (f *FileCache) Get(key string) (interface{}, error) {
	data, err := os.ReadFile(f.filePath(key))
	if err != nil {
		return nil, err
	}
	var val interface{}
	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}
	return val, nil
}

func (f *FileCache) Delete(key string) error {
	return os.Remove(f.filePath(key))
}
