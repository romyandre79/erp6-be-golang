package logger

import (
	"context"
	"erp6-be-golang/core/configs"
	"errors"
)

type Logger interface {
	Info(ctx context.Context, msg string, meta map[string]interface{})
	Error(ctx context.Context, msg string, meta map[string]interface{})
}

// InitLogger
func InitLogger(cfg *configs.Config) (Logger, error) {
	err := errors.New("unsupported, mode: file, remote, db")
	switch cfg.LogMode {
	case "file":
		return NewFileLogger(), nil
	case "remote":
		return NewRemoteLogger(), nil
	default:
		return nil, err
	}
}
