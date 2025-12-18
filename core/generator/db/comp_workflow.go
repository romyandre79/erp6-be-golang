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
		

		// Save parent workflow engine state and termination flag
		parentWfEngine := c.Locals("wfEngine")
		parentTerminated := c.Locals("flowTerminated")

		// Execute the sub-workflow (this will overwrite "wfEngine" in Locals)
		err := ExecuteFlow(c, db, wfName, search, flowParams)

		// Capture sub-workflow results
		var childWfEngine []WorkflowEngine
		if c.Locals("wfEngine") != nil {
			childWfEngine = c.Locals("wfEngine").([]WorkflowEngine)
		}

		// Restore parent workflow engine state and termination flag
		// This ensures sub-workflow "End" node doesn't kill the parent workflow
		c.Locals("wfEngine", parentWfEngine)
		c.Locals("flowTerminated", parentTerminated)

		if err != nil {
			return err
		}

		// Propagate variables from child workflow to parent
		// We collect all ResultNode maps from the child workflow
		mergedResults := make(map[string]interface{})
		for _, step := range childWfEngine {
			if step.ResultNode != nil {
				if resMap, ok := step.ResultNode.(map[string]interface{}); ok {
					for k, v := range resMap {
						mergedResults[k] = v
						fmt.Printf("[Workflow] Propagating variable $%s = %v from sub-flow\n", k, v)
					}
				}
			}
		}
		
		// Note: The InternalFlow (caller) will append the result of THIS component to wfEngine.
		// However, InternalFlow uses the RETURN value of handleWorkflow? 
		// No, handleWorkflow returns 'error'. 
		// InternalFlow checks 'wfEngine' for the LAST entry.
		
		// If we want InternalFlow to register our merged results, we must append it ourselves 
		// OR let InternalFlow handle it.
		// But InternalFlow logic is: 
		// "if len(wfEngine) > 0 && wfEngine[len-1].ComponentName == component.Name"
		// Since we restored parentWfEngine, the last entry is NOT us (it's the previous node).
		
		// So we MUST append our result node manually to parentWfEngine.
		
		finalWfEngine := c.Locals("wfEngine").([]WorkflowEngine)
		finalWfEngine = append(finalWfEngine, WorkflowEngine{
			ResultNode: mergedResults,
			Success:    true,
		})
		c.Locals("wfEngine", finalWfEngine)
		
		return nil
	} else {
		fmt.Printf("[Workflow] Workflow disabled, skipping execution\n")
		return nil
	}
}
