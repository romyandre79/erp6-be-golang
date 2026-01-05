package generator

import (
	"fmt"
	"strings"

	"erp6-be-golang/core/helpers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("table", func(ctx *WorkflowContext) error {
		return handleTable(ctx)
	})
}

func handleTable(ctx *WorkflowContext) error {
	c := ctx.FiberCtx
	params := ctx.Params
	db := ctx.DB
	
	var (
		param           string
		tablename       string
		method          string
		enable          = true
		skipTransaction = false
	)

	for _, p := range params {
		switch p.InputName {
		case "tableparam":
			param = strings.TrimSpace(p.CompValue)
		case "table":
			tablename = strings.TrimSpace(p.CompValue)
		case "method":
			method = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "enabletable":
			if p.CompValue == "false" {
				enable = false
			}
		case "skip_transaction", "skiptransaction":
			if strings.ToLower(p.CompValue) == "true" {
				skipTransaction = true
			}
		}
	}

	// ... (parameter parsing logic remains same, skipping lines 46-126 in replacement for brevity if tools allow, but here I must match context)
	// Actually, I can just replace the top block and the DB selection block. 
	// But `replace_file_content` needs contiguous block. 
	// I'll replace the top var block first.
	// Wait, I need to do this in one go or multiple steps. 
	// The file size is small enough.
	
	// Let's do the VAR block first.


	// parse params respecting quotes
	var listOldParam []string
	var currentParam strings.Builder
	inQuote := false
	for _, r := range param {
		if r == '\'' {
			inQuote = !inQuote
			currentParam.WriteRune(r)
		} else if r == ',' && !inQuote {
			listOldParam = append(listOldParam, strings.TrimSpace(currentParam.String()))
			currentParam.Reset()
		} else {
			currentParam.WriteRune(r)
		}
	}
	if currentParam.Len() > 0 {
		listOldParam = append(listOldParam, strings.TrimSpace(currentParam.String()))
	}

	postData := map[string]string{}

	// ambil semua data POST
	form, _ := c.MultipartForm()
	if form != nil {
		for key, val := range form.Value {
			postData[key] = val[0]
		}
	}

	// build parameter baru
	newParam := map[string]interface{}{}
	
	// Debug: Print available Extras
	if ctx.Extras != nil {
		fmt.Println("[CompTable] Available Extras:")
		for k, v := range ctx.Extras {
			fmt.Printf("  - %s: %v\n", k, v)
		}
	} else {
		fmt.Println("[CompTable] No Extras available")
	}

	for _, key := range listOldParam {
		if strings.Contains(key, "=") {
			parts := strings.SplitN(key, "=", 2)
			valRaw := strings.TrimSpace(parts[1])

			// Strip surrounding quotes if present
			if len(valRaw) >= 2 && strings.HasPrefix(valRaw, "'") && strings.HasSuffix(valRaw, "'") {
				valRaw = valRaw[1 : len(valRaw)-1]
			}

			if strings.Contains(valRaw, "$") {
				// Use standardized ResolveParam to handle all sources (Extras, Locals, Query, NodeResults)
				val := ResolveParam(c, valRaw)
				fmt.Printf("[CompTable] Resolved variable '%s' to '%s'\n", valRaw, val)
				
				// Keep fallback logic if ResolveParam returns the variable name itself (meaning not found)
				// modifying params only if resolved
				if val != valRaw {
					newParam[parts[0]] = val
				} else {
					// Check if we wanted to force empty string for unresolved vars?
					// Old logic printed "Failed to resolve... using empty string" but seemingly used newParam[parts[0]] = "" implicitly via zero value?
					// Actually old logic: `var val string` (empty) -> if not found -> val remains "" -> `newParam[parts[0]] = val`
					// ResolveParam returns "$varname" if not found.
					// So we should handle that.
					fmt.Printf("[CompTable] Failed to resolve '%s', utilizing default empty string\n", valRaw)
					newParam[parts[0]] = ""
				}
			} else {
				newParam[parts[0]] = valRaw
			}
		} else {
			// Standalone key (e.g. "modulename")
			// Priority 1: Check POST/Form data
			if val, ok := postData[key]; ok {
				newParam[key] = val
			} else {
				// Priority 2: Try to auto-resolve as variable (e.g. $modulename)
				// This fixes the issue where "modulename" became 1 because it wasn't explicitly "$modulename"
				resolved := ResolveParam(c, "$"+key)
				
				// ResolveParam returns the input ("$key") if not found
				if resolved != "$"+key {
					newParam[key] = resolved
					fmt.Printf("[CompTable] Auto-resolved standalone key '%s' to '%s'\n", key, resolved)
				} else {
					// Fallback to 1 if not found anywhere (legacy behavior)
					newParam[key] = 1
				}
			}
		}
	}

	// mulai proses SQL dinamis
	// Determine which DB connection to use (Transaction vs Raw)
	if skipTransaction {
		if c != nil {
			if rawDB, ok := c.Locals("db").(*gorm.DB); ok {
				db = rawDB
				fmt.Println("[CompTable] bypassing transaction (skip_transaction=true)")
			} else {
				fmt.Println("[CompTable] Warning: skip_transaction=true but raw DB not found in Locals, using default (transactional)")
			}
		}
	}

	// mulai proses SQL dinamis
	// Insert
	if strings.HasPrefix(strings.ToLower(method), "insert") {
		// Use raw DB (no transaction) if skipTransaction is requested
		useDB := db
		if skipTransaction {
			if rawDB, err := GetRawDBConnection(ctx.FiberCtx); err == nil {
				useDB = rawDB
				fmt.Println("Using Raw DB connection for INSERT (skipping workflow transaction)")
			}
		}

		result := useDB.Table(tablename).Create(newParam)
		if result.Error != nil {
			return result.Error
		}
	}

	// mulai proses SQL dinamis
	switch method {
	case "insert":
		if enable {
			fmt.Printf("[CompTable] Executing INSERT on table '%s' (SkipTransaction: %v)\n", tablename, skipTransaction)
			fmt.Printf("[CompTable] Payload: %+v\n", newParam)

			result := db.Table(tablename).Create(newParam)
			if result.Error != nil {
				fmt.Printf("[CompTable] INSERT FAILED: %v\n", result.Error)
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA CREATE", "TABLE "+tablename)
				return result.Error
			}
			fmt.Printf("[CompTable] INSERT SUCCESS. RowsAffected: %d\n", result.RowsAffected)


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
			postData["lastid"] = fmt.Sprint(lastID)
			
			// Construct response data
			responseData := make(map[string]interface{})
			for k, v := range newParam {
				responseData[k] = v
			}
			responseData["lastid"] = lastID
			
			responseData["lastid"] = lastID
			
			helpers.SuccessResponse(c, "DATA SAVED", responseData)

		} else {
			result := db.Table(tablename).Session(&gorm.Session{DryRun: true}).Create(newParam)
			rawQuery := result.Statement.SQL.String()
			rawVars := result.Statement.Vars
			helpers.SuccessResponse(c, "DATA SAVED", map[string]interface{}{
				"sql":  rawQuery,
				"vars": rawVars,
			})
		}

	case "update":
		idField := listOldParam[0]
		idValue := postData[idField]
		delete(newParam, idField)

		if enable {
			result := db.Table(tablename).Where(fmt.Sprintf("%s = ?", idField), idValue).Updates(newParam)
			if result.Error != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA UPDATE", "TABLE "+tablename)
				return result.Error
			}
			helpers.SuccessResponse(c, "DATA SAVED", "")
		} else {
			result := db.Table(tablename).Where(fmt.Sprintf("%s = ?", idField), idValue).Session(&gorm.Session{DryRun: true}).Updates(newParam)
			rawQuery := result.Statement.SQL.String()
			rawVars := result.Statement.Vars
			helpers.SuccessResponse(c, "DATA SAVED", map[string]interface{}{
				"sql":  rawQuery,
				"vars": rawVars,
			})
		}

	case "purge":
		idField := listOldParam[0]
		idValue := postData[idField]
		sqlment := fmt.Sprintf("delete from %s where %s = %s", tablename, idField, idValue)

		if enable {
			result := db.Exec(sqlment)
			if result.Error != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA PURGE", "TABLE "+tablename)
				return result.Error
			}
			helpers.SuccessResponse(c, "DATA SAVED", "")
		} else {
			helpers.FailResponse(c, 401, "INVALID_DATA SAVED", sqlment)
		}

	default:
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA UPDATE", "TABLE "+tablename)
	}
	return nil
}
