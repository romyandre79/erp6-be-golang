// internal/core/logger/file_logger.go
package logger

import (
    "context"
    "log"
    "os"
)

type FileLogger struct {
    logger *log.Logger
}

func NewFileLogger() *FileLogger {
    file, _ := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    return &FileLogger{
        logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
    }
}

func (f *FileLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
    f.logger.Printf("[INFO] %s %v\n", msg, fields)
}

func (f *FileLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
    f.logger.Printf("[ERROR] %s %v\n", msg, fields)
}

func (f *FileLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
    f.logger.Printf("[WARN] %s %v\n", msg, fields)
}
