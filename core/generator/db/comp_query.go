package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"erp6-be-golang/core/helpers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("query", func(ctx *WorkflowContext) error {
		return handleQuery(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

func handleQuery(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		datatype  string
		query     string
		parameter string
	)

	for _, p := range params {
		switch p.InputName {
		case "datatype":
			datatype = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "query":
			query = strings.TrimSpace(p.CompValue)
		case "parameter":
			parameter = strings.TrimSpace(p.CompValue)
		}
	}

	resultStat := make(map[string]interface{})

	// Parse parameters if needed (expecting JSON array for SQL params)
	var sqlParams []interface{}
	if parameter != "" {
		// Try to parse as JSON array
		if err := json.Unmarshal([]byte(parameter), &sqlParams); err != nil {
			// If not JSON array, check if it's a comma-separated list
			if strings.Contains(parameter, ",") {
				parts := strings.Split(parameter, ",")
				sqlParams = make([]interface{}, len(parts))
				for i, p := range parts {
					sqlParams[i] = strings.TrimSpace(p)
				}
			} else {
				// Treat as single value
				sqlParams = []interface{}{parameter}
			}
		}

		// Resolve variables in parameters
		for i, param := range sqlParams {
			if strParam, ok := param.(string); ok {
				resolved := strParam
				if strings.Contains(strParam, "$") {
					resolved = ResolveParam(c, strParam)
					fmt.Printf("[Query Debug] Parameter '%s' resolved to '%s'\n", strParam, resolved)
				}

				// Auto-fix number format for ID/EU currency (e.g. 16.636,40 -> 16636.40, or 11,53 -> 11.53)
				// Condition: Contains comma, check if it's likely a decimal separator (appears after any dots)
				// Trim space first to ensure clean check
				resolved = strings.TrimSpace(resolved)
				
				if strings.Contains(resolved, ",") {
					isNumeric := true
					for _, r := range resolved {
						if (r < '0' || r > '9') && r != '.' && r != ',' {
							isNumeric = false
							fmt.Printf("[Query Debug] Normalization skipped for '%s': Found invalid char '%c'\n", resolved, r)
							break
						}
					}

					if isNumeric {
						lastComma := strings.LastIndex(resolved, ",")
						lastDot := strings.LastIndex(resolved, ".")
						
						// If comma is the last separator (or the only separator)
						// e.g. 1.234,56 (comma > dot)
						// e.g. 11,53 (dot is -1, comma > -1)
						if lastComma > lastDot {
							original := resolved
							// Remove all dots (thousand separators)
							resolved = strings.ReplaceAll(resolved, ".", "")
							// Replace comma with dot (decimal separator)
							resolved = strings.ReplaceAll(resolved, ",", ".")
							fmt.Printf("[Query Debug] Auto-normalized number: '%s' -> '%s'\n", original, resolved)
						} else {
                            fmt.Printf("[Query Debug] Normalization skipped for '%s': Comma found but Dot is later (likely US format 1,234.56)\n", resolved)
                        }
					}
				}
				
				sqlParams[i] = resolved
			}
		}
	}

	fmt.Printf("[Query Debug] Executing Query: %s with Params: %v\n", query, sqlParams)
	fmt.Printf("[Query Debug] Component Datatype: %s\n", datatype)

	switch strings.ToLower(datatype) {
	case "node_result":
		// ... existing node_result logic ...
		// Just return the parameter (previous result) as the result of this node
		// If parameter was JSON, we might want to return the parsed object?
		// "accept output from previous result"
		// If the input was a JSON string, unmarshaling it above gives us the generic structure.
		// If the user wants the raw string, we can give that too.
		// Usually we want the structured data if possible.
		if len(sqlParams) > 0 {
			// If it was a list of params, return that list?
			// Or if it was a complex object?
			// Let's try to unmarshal as arbitrary interface for node_result
			var genericData interface{}
			if err := json.Unmarshal([]byte(parameter), &genericData); err == nil {
				resultStat["data"] = genericData
			} else {
				resultStat["data"] = parameter
			}
		} else {
			resultStat["data"] = parameter
		}

		wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
		c.Locals("wfEngine", wfEngine)

		helpers.SuccessResponse(c, "DATA RETRIEVED", resultStat)

	case "get":
		if query == "" {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA", "EMPTY QUERY")
			return nil
		}
		
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(query)), "insert") || 
		   strings.HasPrefix(strings.ToLower(strings.TrimSpace(query)), "update") ||
		   strings.HasPrefix(strings.ToLower(strings.TrimSpace(query)), "delete") {
			fmt.Printf("[Query Debug] WARNING: Running INSERT/UPDATE/DELETE with Method 'GET'. This may loop incorrectly or fail to commit depending on driver.\n")
		}

		var rows []map[string]interface{}
		if err := db.Raw(query, sqlParams...).Scan(&rows).Error; err != nil {
			helpers.FailResponse(c, fiber.StatusInternalServerError, "QUERY ERROR", err.Error())
			return err
		}

		resultStat["data"] = rows

		// Add formatted message for SendMessage component
		resultStat["message"] = formatDataAsTable(rows)

		// Setup paging meta if relevant? (Not implemented for raw query yet)

		wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
		c.Locals("wfEngine", wfEngine)

		helpers.SuccessResponse(c, "DATA RETRIEVED", resultStat)

	case "post":
		if query == "" {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA", "EMPTY QUERY")
			return nil
		}

		result := db.Exec(query, sqlParams...)
		if result.Error != nil {
			helpers.FailResponse(c, fiber.StatusInternalServerError, "EXEC ERROR", result.Error.Error())
			return result.Error
		}
		
		fmt.Printf("[Query Debug] Rows Affected: %d\n", result.RowsAffected)

		// Try to get Last Insert ID if possible (Driver dependent)
		var lastID int64
		// Simple approach: only if RowsAffected > 0 and it's an insert?
		// Let's just emulate comp_table best effort
		driver := GetDatabaseDriver(db) // Assuming this helper exists, check comp_table

		// Copied from comp_table logic
		// Only attempt to get Last Insert ID if rows were actually affected
		if result.RowsAffected > 0 && strings.HasPrefix(strings.ToLower(query), "insert") {
			switch driver {
			case "postgres":
				db.Raw("SELECT LASTVAL()").Scan(&lastID)
			case "sqlserver":
				db.Raw("SELECT CAST(COALESCE(SCOPE_IDENTITY(), 0) AS BIGINT)").Scan(&lastID)
			case "sqlite", "sqlite3":
				db.Raw("SELECT last_insert_rowid()").Scan(&lastID)
			case "oracle":
				lastID = 0
			case "mysql", "mariadb":
				db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
			default:
				// Try MySQL syntax as default
				db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
			}
			resultStat["lastid"] = lastID
		}

		resultStat["rows_affected"] = result.RowsAffected

		wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
		c.Locals("wfEngine", wfEngine)

		helpers.SuccessResponse(c, "DATA SAVED", resultStat)

	default:
		helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID DATATYPE", datatype)
	}

	return nil
}
