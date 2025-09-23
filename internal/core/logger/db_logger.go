package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DBLogger implements Logger interface
type DBLogger struct {
	db *gorm.DB
}

// LogModel schema untuk menyimpan log ke DB
type LogModel struct {
	ID        uint      `gorm:"primaryKey"`
	Level     string    `gorm:"size:20"`
	Message   string    `gorm:"size:500"`
	Metadata  string    `gorm:"type:text"`
	CreatedAt time.Time
}

// NewDBLogger membuat instance DBLogger
func NewDBLogger(db *gorm.DB) *DBLogger {
	// Auto migrate log table
	db.AutoMigrate(&LogModel{})
	return &DBLogger{db: db}
}

func (l *DBLogger) Info(ctx context.Context, msg string, meta map[string]interface{}) {
	l.saveLog("INFO", msg, meta)
}

func (l *DBLogger) Error(ctx context.Context, msg string, meta map[string]interface{}) {
	l.saveLog("ERROR", msg, meta)
}

func (l *DBLogger) saveLog(level, msg string, meta map[string]interface{}) {
	log := LogModel{
		Level:     level,
		Message:   msg,
		Metadata:  fmt.Sprintf("%v", meta),
		CreatedAt: time.Now(),
	}
	l.db.Create(&log)
}
