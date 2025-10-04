package storage

import (
	"erp6-be-golang/core/configs"
	"mime/multipart"
)

type Storage interface {
	Name() string
	Upload(file *multipart.FileHeader, path string) (string, error)
	Delete(path string) error
}

var registry = map[string]Storage{}

func Register(s Storage) {
	registry[s.Name()] = s
}

func Get() (Storage, bool) {
	name := configs.ConfigApps.StorageType
	if name != "" {
		s, ok := registry[name]
		return s, ok
	} else {
		return nil, false
	}
}
