package generator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func init() {
	RegisterComponent("webhook", func(ctx *WorkflowContext) error {
		return handleWebhook(ctx.FiberCtx)
	})
}

func handleWebhook(c *fiber.Ctx) error {
	fmt.Printf("[Webhook] Received request: %s %s\n", c.Method(), c.OriginalURL())
	resultStat := make(map[string]interface{})

	// Capture Body
	var body map[string]interface{}
	if err := c.BodyParser(&body); err == nil {
		resultStat["body"] = body
	} else {
		// If body is not JSON, try to capture as simple string or form
		if len(c.Body()) > 0 {
			resultStat["raw_body"] = string(c.Body())
		}
	}

	// Capture Query Params
	queries := c.Queries()
	if len(queries) > 0 {
		resultStat["query"] = queries
	}

	// Capture Headers
	headers := c.GetReqHeaders()
	if len(headers) > 0 {
		resultStat["headers"] = headers
	}

	// Capture Form Data
	if form, err := c.MultipartForm(); err == nil {
		resultStat["form"] = form.Value
	}

	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{
		ComponentName: "webhook",
		DataInputNode: "",
		ResultNode:    resultStat,
		Success:       true,
	})
	c.Locals("wfEngine", wfEngine)

	return nil
}
