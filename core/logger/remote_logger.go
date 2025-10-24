package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// RemoteLogger kirim log ke server eksternal
type RemoteLogger struct {
	endpoint string
	client   *http.Client
}

type remotePayload struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewRemoteLogger bikin instance RemoteLogger
func NewRemoteLogger() *RemoteLogger {
	endpoint := os.Getenv("LOG_REMOTE_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:9000/logs" // default
	}
	return &RemoteLogger{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (l *RemoteLogger) Info(ctx context.Context, msg string, meta map[string]interface{}) {
	l.send("INFO", msg, meta)
}

func (l *RemoteLogger) Error(ctx context.Context, msg string, meta map[string]interface{}) {
	l.send("ERROR", msg, meta)
}

func (l *RemoteLogger) send(level, msg string, meta map[string]interface{}) {
	payload := remotePayload{
		Level:     level,
		Message:   msg,
		Metadata:  meta,
		Timestamp: time.Now(),
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", l.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	_, _ = l.client.Do(req)
}
