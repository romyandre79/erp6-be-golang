package generator

import (
	"fmt"
	"strings"
)

func init() {
	fmt.Println("[Init] Registering component 'authdata' and 'auth_data'")
	// Register both variants to be safe
	RegisterComponent("authdata", func(ctx *WorkflowContext) error {
		return handleAuthData(ctx)
	})
	RegisterComponent("auth_data", func(ctx *WorkflowContext) error {
		return handleAuthData(ctx)
	})
}

// handleAuthData filters data based on user group rules
// Expected input parameter values format: "groupname:variable_name:value"
// Example: 
// - "admin:filter_type:all"
// - "sales:filter_region:1,2,3,4,5" (Comma separated values are preserved as a single string)
// Logic:
// 1. Get current user's groups
// 2. Iterate through all node parameters
// 3. For each parameter, parse "group:var:val"
// 4. If user belongs to 'group', set $var = val in the workflow engine results
// 5. If at least one rule matches, ctx.DecisionResult = true
func handleAuthData(ctx *WorkflowContext) error {
	// 1. Get User ID
	fmt.Println("[AuthData] Handler START")
	
	userID, ok := ctx.FiberCtx.Locals("userid").(int)
	if !ok || userID == 0 {
		fmt.Println("[AuthData] No user ID found in context locals")
		ctx.DecisionResult = false
		return nil
	}

	// 2. Get User's Groups
	// We need to fetch groups to filter the rules correctly.
	var userGroups []string
	err := ctx.DB.Table("usergroup ug").
		Select("ga.groupname").
		Joins("JOIN groupaccess ga ON ga.groupaccessid = ug.groupaccessid").
		Where("ug.useraccessid = ?", userID).
		Pluck("ga.groupname", &userGroups).Error

	if err != nil {
		fmt.Printf("[AuthData] DB Error fetching user groups: %v\n", err)
		// Don't fail the flow, just log error and proceed with empty groups (no matches)
		// return helpers.FailResponse(ctx.FiberCtx, fiber.StatusInternalServerError, "DB_ERROR", fmt.Sprintf("Auth Data check failed: %v", err))
	}

	// Simplify group check with a map
	userGroupMap := make(map[string]bool)
	for _, g := range userGroups {
		userGroupMap[strings.ToLower(strings.TrimSpace(g))] = true
	}
	fmt.Printf("[AuthData] User Groups: %v\n", userGroups)

	// 3. Process Rules
	authorized := false
	resultMap := make(map[string]interface{})

	for i, param := range ctx.Params {
		fmt.Printf("[AuthData] Processing Param [%d]: Name='%s', Value='%s'\n", i, param.InputName, param.CompValue)
		
		fullRuleString := ResolveParam(ctx.FiberCtx, param.CompValue)
		if strings.TrimSpace(fullRuleString) == "" {
			continue
		}

		// Split by comma to handle multiple rules OR comma-separated values
		tokens := strings.Split(fullRuleString, ",")
		
		var currentGroup string
		var currentVar string
		var currentValueBuilder strings.Builder
		
		// Helper to commit the current rule being built
		commitRule := func() {
			if currentGroup != "" && currentVar != "" {
				finalValue := strings.TrimSpace(currentValueBuilder.String())
				
				// Check group (case-insensitive)
				if userGroupMap[strings.ToLower(currentGroup)] {
					fmt.Printf("[AuthData] Match: User in group '%s', setting $%s = '%s'\n", currentGroup, currentVar, finalValue)
					resultMap[currentVar] = finalValue
					authorized = true
				} else {
					fmt.Printf("[AuthData] Skip: User NOT in group '%s' (Value: '%s')\n", currentGroup, finalValue)
				}
			}
		}

		for _, token := range tokens {
			// Check for "Group:Var:Value" pattern (approx 2 colons)
			if strings.Count(token, ":") >= 2 {
				// Commit previous rule if exists
				commitRule()
				
				// Start new rule
				parts := strings.SplitN(token, ":", 3) // Split into max 3 parts: Group, Var, Rest
				currentGroup = strings.TrimSpace(parts[0])
				currentVar = strings.TrimSpace(parts[1])
				
				currentValueBuilder.Reset()
				currentValueBuilder.WriteString(parts[2]) // Start value
			} else {
				// Continuation of previous value
				if currentValueBuilder.Len() > 0 {
					currentValueBuilder.WriteString(",")
					currentValueBuilder.WriteString(token)
				}
			}
		}
		// Commit the last rule
		commitRule()
	}

	// 4. Register results and decision
	if authorized {
		// Append to wfEngine so subsequent nodes can use $variable
		wfEngine := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		
		wfEngine = append(wfEngine, WorkflowEngine{
			ResultNode: resultMap,
			Success:    true,
		})
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
	}

	ctx.DecisionResult = authorized
	if authorized {
		fmt.Printf("[AuthData] User %d Authorized. Variables set: %v\n", userID, resultMap)
	} else {
		fmt.Printf("[AuthData] User %d Rejected (No matching group rules)\n", userID)
	}

	return nil
}
