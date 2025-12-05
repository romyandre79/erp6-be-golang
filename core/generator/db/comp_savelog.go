package generator

import (
	"context"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/logger"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("savelog", func(ctx *WorkflowContext) error {
		return handleSaveLog(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

// handleSaveLog saves log messages based on component parameters
// Supported parameters:
// - logtype: logging destination (file, remote, db) - default: file
// - logdatatype: type of log data (apps, other, db) - categorizes the log
// - logcontent: the log message content (required)
func handleSaveLog(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		logType     = "file"
		logDataType = "apps"
		logContent  string
		metadata    = make(map[string]interface{})
	)

	// Extract parameters from workflow
	for _, p := range params {
		switch strings.ToLower(p.InputName) {
		case "logtype":
			logType = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "logdatatype":
			logDataType = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "logcontent":
			logContent = strings.TrimSpace(p.CompValue)
		}
	}

	// Validate required parameters
	if logContent == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_PARAMETER", "logcontent is required")
	}

	// Validate logtype
	if logType != "file" && logType != "remote" && logType != "db" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_LOG_TYPE", "logtype must be file, remote, or db")
	}

	// Validate logdatatype
	if logDataType != "apps" && logDataType != "other" && logDataType != "db" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_LOG_DATATYPE", "logdatatype must be apps, other, or db")
	}

	// Add user context to metadata if available
	if username, ok := c.Locals("username").(string); ok && username != "" {
		metadata["username"] = username
	}
	if userID, ok := c.Locals("user_id").(int); ok && userID > 0 {
		metadata["user_id"] = userID
	}

	// Add request context information
	metadata["datatype"] = logDataType
	metadata["logtype"] = logType
	metadata["ip"] = c.IP()
	metadata["method"] = c.Method()
	metadata["path"] = c.Path()
	metadata["user_agent"] = c.Get("User-Agent")

	// Add request ID if available
	if reqID := c.Get("X-Request-ID"); reqID != "" {
		metadata["request_id"] = reqID
	}

	ctx := context.Background()
	var logInstance logger.Logger

	// Initialize logger based on logtype parameter
	switch logType {
	case "file":
		logInstance = logger.NewFileLogger()
	case "remote":
		logInstance = logger.NewRemoteLogger()
	case "db":
		logInstance = logger.NewDBLogger(db)
	default:
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "INVALID_LOG_TYPE",
			fmt.Sprintf("unsupported logtype: %s. Supported types: file, remote, db", logType))
	}

	// Save log as INFO level (you can extend this to support level parameter if needed)
	logInstance.Info(ctx, logContent, metadata)
	return nil
}
