package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
	BaseURL  string
}

func (l *LocalStorage) Name() string {
	return "local"
}

func (l *LocalStorage) Upload(file *multipart.FileHeader, path string) (string, error) {
	dstPath := filepath.Join(l.BasePath, path)
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return "", err
	}

	dst := filepath.Join(dstPath, file.Filename)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", l.BaseURL, path, file.Filename), nil
}

func (l *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(l.BasePath, path))
}

func init() {
	Register(&LocalStorage{
		BasePath: "./uploads",
		BaseURL:  "http://localhost:3000/uploads",
	})
}
