package logger

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Logger interface {
	Info(ctx context.Context, msg string, meta map[string]interface{})
	Error(ctx context.Context, msg string, meta map[string]interface{})
}

// NewLogger pilih logger sesuai mode
func NewLogger(mode string) Logger {
	switch mode {
	case "db":
		// Default: SQLite local untuk demo
		db, err := gorm.Open(sqlite.Open("logs.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database for logger")
		}
		return NewDBLogger(db)
	case "remote":
		return NewRemoteLogger()
	default:
		return NewFileLogger() // fallback ke file logger
	}
}
