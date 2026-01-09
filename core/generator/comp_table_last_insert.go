package generator 

import (
	"erp6-be-golang/core/helpers"
	"strings"
)

func init() {
	RegisterComponent("TableLastInsert", func(ctx *WorkflowContext) error {
		return handleTableLastInsert(ctx)
	})
}

func handleTableLastInsert(ctx *WorkflowContext) error {
	c := ctx.FiberCtx
	db := ctx.DB
	params := ctx.Params

	var (
		parambox           string
		enablebox          = true
	)

	for _, p := range params {
		switch p.InputName {
		case "parambox":
			parambox = strings.TrimSpace(p.CompValue)
		case "enablebox":
			if p.CompValue == "false" {
				enablebox = false
			}
		}
	}

	var lastID int64
	driver := GetDatabaseDriver(db)
	switch driver {
		case "postgres":
			db.Raw("SELECT LASTVAL()").Scan(&lastID)
		case "sqlserver":
			db.Raw("SELECT CAST(COALESCE(SCOPE_IDENTITY(), 0) AS BIGINT)").Scan(&lastID)
		case "sqlite", "sqlite3":
			db.Raw("SELECT last_insert_rowid()").Scan(&lastID)
		case "oracle":
			// Oracle with map insert makes capturing ID difficult without RETURNING clause support in finding the identity sequence
			// Use 0 or handle specifically if needed.
			lastID = 0
		default: // mysql, mariadb
			db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
	}
	responseData := make(map[string]interface{})
	ctx.Extras["lastID"] = lastID
	responseData["lastID"] = lastID

	if parambox != "" {
		ctx.Extras[parambox] = lastID
		responseData[parambox] = lastID
	}

	if enablebox == true {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: responseData})
	c.Locals("wfEngine", wfEngine)
	
	helpers.SuccessResponse(c, "DATA SAVED", responseData)
	}
	return nil
}