package generator

import (
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
	enable := true
	for _, v := range params {
		switch v.InputName {
		case "workflowname":
			wfName = strings.TrimSpace(v.CompValue)
		//case "workflowparameter":
		//	wfParameter = strings.TrimSpace(v.CompValue)
		case "enableworkflow":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}
	if enable {
		return ExecuteFlow(c, db, wfName, search)
	} else {
		return nil
	}
}
