package generator

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	RegisterComponent("decision", func(ctx *WorkflowContext) error {
		return handleDecision(ctx)
	})
}

// handleDecision evaluates a condition and sets the decision result in the context
// Supported parameters:
// - decisionok: the condition to evaluate (format: "key=value" or "key=empty")
// - decisionparamtype: source of the value (post, get, node result) - default: post
// - enabledecision: whether to enable decision validation (true/false) - default: true
//
// Decision logic:
// - Compares a value from the specified source against an expected value
// - Sets ctx.DecisionResult to true if condition matches, false otherwise
// - If expected value is "empty", checks if the value is empty string
func handleDecision(ctx *WorkflowContext) error {
	decisionParamType := "post" // default
	contentDecision := ""
	enableDecision := true

	var wfEngine = ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)

	// Parse parameters
	for _, v := range ctx.Params {
		switch v.InputName {
		case "decisionok":
			contentDecision = v.CompValue
		case "decisionparamtype":
			decisionParamType = strings.ToLower(v.CompValue)
		case "enabledecision":
			if strings.ToLower(v.CompValue) == "false" {
				enableDecision = false
			}
		}
	}

	// Parse decision content
	operator := "="
	if strings.Contains(contentDecision, "!=") {
		operator = "!="
	} else if !strings.Contains(contentDecision, "=") {
		if enableDecision {
			return helpers.FailResponse(ctx.FiberCtx, fiber.StatusBadRequest,
				"INVALID_DECISION_FORMAT", "Decision format must be 'key=value' or 'key!=value'")
		}
		ctx.DecisionResult = false
		return nil
	}

	data := strings.SplitN(contentDecision, operator, 2)
	key := strings.TrimSpace(data[0])
	expected := strings.TrimSpace(data[1])

	// Get actual value based on parameter type
	var actualValue string

	switch decisionParamType {
	case "post":
		actualValue = ctx.FiberCtx.FormValue(key)

	case "get":
		actualValue = ctx.FiberCtx.Query(key)

	case "node result", "noderesult":
		if len(wfEngine) > 0 {
			// Get the last node result
			lastResult := wfEngine[len(wfEngine)-1].ResultNode

			// Try as string first
			if resultStr, ok := lastResult.(string); ok {
				actualValue = resultStr
			} else if resultMap, ok := lastResult.(map[string]interface{}); ok {
				// Try to extract the key from the map
				if val, exists := resultMap[key]; exists {
					// Convert value to string
					if strVal, ok := val.(string); ok {
						actualValue = strVal
					} else { // Handle nil or other types
						if val == nil {
							actualValue = "" // explicit empty for nil
						} else {
							actualValue = fmt.Sprintf("%v", val)
						}
					}
				}
			}
		}

	default:
		if enableDecision {
			return helpers.FailResponse(ctx.FiberCtx, fiber.StatusBadRequest,
				"INVALID_PARAM_TYPE",
				fmt.Sprintf("Unknown decisionparamtype: %s. Valid options: post, get, node result", decisionParamType))
		}
		ctx.DecisionResult = false
		return nil
	}

	// Debug logging
	fmt.Printf("[Decision] Condition: %s, Key: %s, Expected: %s, Actual: '%s'\n", 
		contentDecision, key, expected, actualValue)

	// Evaluate decision
	match := false
	if strings.ToLower(expected) == "empty" {
		match = (actualValue == "")
	} else {
		match = (actualValue == expected)
	}

	if operator == "!=" {
		ctx.DecisionResult = !match
	} else {
		ctx.DecisionResult = match
	}

	fmt.Printf("[Decision] Result: %v (match=%v, operator=%s)\n", ctx.DecisionResult, match, operator)
	return nil
}
