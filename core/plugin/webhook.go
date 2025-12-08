package plugin

import (
	"bytes"
	"encoding/json"
	"erp6-be-golang/core/configs"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
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

func WebhookHandler(c *fiber.Ctx) error {
	source := c.Params("source")
	if source == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Source is required"})
	}

	// Validate source against external components
	// We assume external components are folders in ExternalComponentPath
	componentName := strings.ToLower(source)
	// Some sources might map to a component name differently, but standard convention component_name == source

	// We check if the binary exists
	// Convention: ExternalComponentPath/{componentName}/{componentName}.exe (windows) or just {componentName}
	// Let's assume binary name matches component name

	extPath := configs.ConfigApps.ExternalComponentPath
	if extPath == "" {
		extPath = "./temp_external_components"
	}

	binaryPath := filepath.Join(extPath, componentName, componentName+".exe") // Assuming Windows environment as per user instructions
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Try without .exe
		binaryPath = filepath.Join(extPath, componentName, componentName)
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Component %s not found", componentName)})
		}
	}

	// Prepare Input for Component
	// We need to pass the request details to the component so it can "recognize" and "handle" it.
	// Common input structure for components seems to be `{"params": [...]}`
	// We will map request details into params.

	body := c.Body()
	// headers := c.GetReqHeaders() // Fiber v2
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
		// Try to capture stderr if available or return generic error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to execute component: %v", err)})
	}

	// Parse Output
	var output ComponentOutput
	if err := json.Unmarshal(out.Bytes(), &output); err != nil {
		// If output is not JSON, return as string
		return c.Status(http.StatusOK).Send(out.Bytes())
	}

	if output.Error != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": output.Error})
	}

	return c.Status(http.StatusOK).JSON(output.Result)
}
