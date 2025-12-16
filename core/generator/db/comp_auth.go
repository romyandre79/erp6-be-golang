package generator

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	RegisterComponent("auth", func(ctx *WorkflowContext) error {
		return handleAuth(ctx)
	})
}

// handleAuth checks if the user is authorized based on user groups
// Supported parameters:
// - groups: comma separated list of allowed group names (e.g. "admin,mastermaterial")
// Result:
// - Sets ctx.DecisionResult to true if authorized, false otherwise
func handleAuth(ctx *WorkflowContext) error {
	var allowedGroups []string

	// Parse parameters
	for _, v := range ctx.Params {
		// Accept "groups" or "usergroup" as parameter name
		if strings.EqualFold(v.InputName, "groups") || strings.EqualFold(v.InputName, "usergroup") {
			// Resolve parameter value (handle variables like $param)
			resolvedValue := ResolveParam(ctx.FiberCtx, v.CompValue)
			
			parts := strings.Split(resolvedValue, ",")
			for _, p := range parts {
				if trimmed := strings.TrimSpace(p); trimmed != "" {
					allowedGroups = append(allowedGroups, trimmed)
				}
			}
		}
	}

	// Default to deny if no groups specified
	if len(allowedGroups) == 0 {
		fmt.Println("[Auth] No groups specified in parameters, denying access")
		ctx.DecisionResult = false
		return nil
	}

	// Get User ID from context
	userID, ok := ctx.FiberCtx.Locals("userid").(int)
	if !ok || userID == 0 {
		// Try to look for user_id in form or previous node results if not in locals
		// But usually it should be in locals for authenticated requests
		fmt.Println("[Auth] No user ID found in context locals, denying access")
		ctx.DecisionResult = false
		return nil
	}

	// Check DB if user belongs to any of the allowed groups
	var count int64
	// Table usergroup links useraccess and groupaccess
	// Table groupaccess contains the group names
	// We need to check if count(*) > 0
	err := ctx.DB.Table("usergroup ug").
		Joins("JOIN groupaccess ga ON ga.groupaccessid = ug.groupaccessid").
		Where("ug.useraccessid = ? AND ga.groupname IN ?", userID, allowedGroups).
		Count(&count).Error

	if err != nil {
		fmt.Printf("[Auth] DB Error checking groups: %v\n", err)
		return helpers.FailResponse(ctx.FiberCtx, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("Auth check failed: %v", err))
	}

	if count > 0 {
		fmt.Printf("[Auth] User %d AUTHORIZED (Matched one of: %v)\n", userID, allowedGroups)
		ctx.DecisionResult = true
	} else {
		fmt.Printf("[Auth] User %d DENIED (Not in: %v)\n", userID, allowedGroups)
		ctx.DecisionResult = false
	}

	return nil
}
