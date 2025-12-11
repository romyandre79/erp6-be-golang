package plugin

import (
	"bytes"
	"encoding/json"
	"erp6-be-golang/core/configs"
	generator "erp6-be-golang/core/generator/db"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WebhookInput struct {
	Params []Param `json:"params"`
}

type Param struct {
	InputName string `json:"inputname"`
	CompValue string `json:"compvalue"`
}

type ComponentOutput struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

// Update usage
func WebhookHandler(c *fiber.Ctx, db *gorm.DB) error {
	source := c.Params("source")
	if source == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Source is required"})
	}

	// 1. Try External Component
	extPath := configs.ConfigApps.ExternalComponentPath
	if extPath == "" {
		extPath = "./temp_external_components"
	}
	componentName := strings.ToLower(source)
	binaryPath := filepath.Join(extPath, componentName, componentName+".exe")
	if _, err := os.Stat(binaryPath); err != nil {
		binaryPath = filepath.Join(extPath, componentName, componentName)
	}

	// If binary exists, execute external logic
	if _, err := os.Stat(binaryPath); err == nil {
		return executeExternalComponent(c, binaryPath)
	}

	// 2. Try Internal Workflow
	// Verify if workflow exists first or just let ExecuteFlow handle it?
	// ExecuteFlow returns error if workflow invalid.
	if err := generator.ExecuteFlow(c, db, source, false); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Source '%s' not found (checked plugins and workflows)", source)})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// If flow executed but didn't write response (body empty), return generic OK
	if len(c.Response().Body()) == 0 {
		return c.JSON(fiber.Map{"status": "ok", "message": "Workflow executed"})
	}
	return nil
}

func executeExternalComponent(c *fiber.Ctx, binaryPath string) error {
	body := c.Body()
	headers := c.GetReqHeaders()
	headerBytes, _ := json.Marshal(headers)

	queryParams := c.Queries()
	queryBytes, _ := json.Marshal(queryParams)

	method := c.Method()

	input := WebhookInput{
		Params: []Param{
			{InputName: "action", CompValue: "handle_webhook"},
			{InputName: "webhook_url", CompValue: c.OriginalURL()},
			{InputName: "method", CompValue: method},
			{InputName: "body", CompValue: string(body)},
			{InputName: "headers", CompValue: string(headerBytes)},
			{InputName: "query", CompValue: string(queryBytes)},
		},
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode input"})
	}

	// Execute Component
	cmd := exec.Command(binaryPath)
	cmd.Stdin = bytes.NewReader(inputBytes)
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to execute component: %v", err)})
	}

	// Parse Output
	var output ComponentOutput
	if err := json.Unmarshal(out.Bytes(), &output); err != nil {
		return c.Status(http.StatusOK).Send(out.Bytes())
	}

	if output.Error != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": output.Error})
	}

	return c.Status(http.StatusOK).JSON(output.Result)
}
