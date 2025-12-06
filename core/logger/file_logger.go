// internal/core/logger/file_logger.go
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type FileLogger struct {
	logger *log.Logger
}

func NewFileLogger() *FileLogger {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
	}

	// Open or create the log file
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		// Fallback to stdout if file cannot be opened
		return &FileLogger{
			logger: log.New(os.Stdout, "[FILE_LOGGER] ", log.LstdFlags|log.Lshortfile),
		}
	}

	return &FileLogger{
		logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
	}
}

// formatFields converts map to readable key=value format
func formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build key=value pairs
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		v := fields[k]
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(pairs, " ")
}

func (f *FileLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	fieldsStr := formatFields(fields)
	if fieldsStr != "" {
		f.logger.Printf("[INFO] %s | %s\n", msg, fieldsStr)
	} else {
		f.logger.Printf("[INFO] %s\n", msg)
	}
}

func (f *FileLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	fieldsStr := formatFields(fields)
	if fieldsStr != "" {
		f.logger.Printf("[ERROR] %s | %s\n", msg, fieldsStr)
	} else {
		f.logger.Printf("[ERROR] %s\n", msg)
	}
}

func (f *FileLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	fieldsStr := formatFields(fields)
	if fieldsStr != "" {
		f.logger.Printf("[WARN] %s | %s\n", msg, fieldsStr)
	} else {
		f.logger.Printf("[WARN] %s\n", msg)
	}
}
