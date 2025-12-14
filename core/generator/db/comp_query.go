package generator

import (
	"encoding/json"
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
			// If not JSON array, treat as single parameter ??
			// Or maybe split by comma if it looks like a list?
			// For now, let's treat non-JSON as single string, or empty slice if invalid
			// But for safety, let's just use it as single scalar if it's not array
			// However, previous components like comp_table used comma-separated keys.
			// But here we likely want direct values from AI component or similar.
			// Let's assume if it's not JSON, it is a single value.
			sqlParams = []interface{}{parameter}
		}
	}

	switch strings.ToLower(datatype) {
	case "node_result":
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

		var rows []map[string]interface{}
		if err := db.Raw(query, sqlParams...).Scan(&rows).Error; err != nil {
			helpers.FailResponse(c, fiber.StatusInternalServerError, "QUERY ERROR", err.Error())
			return err
		}

		resultStat["data"] = rows

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
