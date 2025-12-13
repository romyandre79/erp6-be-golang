package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("workflow", func(ctx *WorkflowContext) error {
		return handleWorkflow(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
}

func handleWorkflow(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	wfName := ""
	wfParameter := ""
	enable := true
	for _, v := range params {
		switch v.InputName {
		case "workflowname":
			wfName = strings.TrimSpace(ResolveParam(c, v.CompValue))
		case "workflowparameter":
			wfParameter = strings.TrimSpace(v.CompValue)
		case "enableworkflow":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}
	if enable {
		// Parse parameters
		flowParams := make(map[string]interface{})
		if wfParameter != "" {
			// Try JSON first
			if err := json.Unmarshal([]byte(ResolveParam(c, wfParameter)), &flowParams); err != nil {
				// Fallback to comma-separated key=value
				// e.g. key1=val1,key2=$var
				pairs := strings.Split(wfParameter, ",")
				for _, pair := range pairs {
					kv := strings.SplitN(pair, "=", 2)
					if len(kv) == 2 {
						key := strings.TrimSpace(kv[0])
						val := strings.TrimSpace(kv[1])
						flowParams[key] = ResolveParam(c, val)
					}
				}
			} else {
				// Resolve values in map if it was JSON
				for k, v := range flowParams {
					if strVal, ok := v.(string); ok {
						flowParams[k] = ResolveParam(c, strVal)
					}
				}
			}
		}

		fmt.Printf("[Workflow] Executing workflow: '%s' with params: %+v\n", wfName, flowParams)
		
		// Set flag to indicate nested workflow execution
		c.Locals("nestedWorkflow", true)
		
		return ExecuteFlow(c, db, wfName, search, flowParams)
	} else {
		fmt.Printf("[Workflow] Workflow disabled, skipping execution\n")
		return nil
	}
}
