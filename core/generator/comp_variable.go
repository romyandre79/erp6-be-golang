package generator

import "strings"

func init() {
	RegisterComponent("Variable", func(ctx *WorkflowContext) error {
		return handleVariable(ctx)
	})
}

func handleVariable(ctx *WorkflowContext) error {
	var (
		sourceVar string
		sourceVal string
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		switch p.InputName {
		case "variable":
			sourceVar = val
		case "value":
			sourceVal = val
		}
	}

	if sourceVar != "" {
		ctx.Extras[sourceVar] = sourceVal
	}

	return nil
}
