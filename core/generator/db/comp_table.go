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
		param     string
		tablename string
		method    string
		enable    = true
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
		}
	}

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
	for _, key := range listOldParam {
		if strings.Contains(key, "=") {
			parts := strings.SplitN(key, "=", 2)
			valRaw := strings.TrimSpace(parts[1])

			// Strip surrounding quotes if present
			if len(valRaw) >= 2 && strings.HasPrefix(valRaw, "'") && strings.HasSuffix(valRaw, "'") {
				valRaw = valRaw[1 : len(valRaw)-1]
			}

			if strings.Contains(valRaw, "$") {
				lookupKey := strings.ReplaceAll(valRaw, "$", "")
				var val string
				
				// Priority 1: Check Extras from previous workflow nodes or conversation state
				if ctx.Extras != nil {
					if extraVal, exists := ctx.Extras[lookupKey]; exists {
						val = fmt.Sprint(extraVal)
					}
				}
				
				// Priority 2: Check Locals (JWT token data)
				if val == "" {
					if lookupKey == "userid" {
						if userID, ok := c.Locals("userid").(int); ok && userID != 0 {
							val = fmt.Sprint(userID)
						}
					} else if lookupKey == "username" {
						if username, ok := c.Locals("username").(string); ok && username != "" {
							val = username
						}
					}
				}
				
				// Priority 3: Check Query and FormValue
				if val == "" {
					val = c.Query(lookupKey)
					if val == "" {
						val = c.FormValue(lookupKey)
					}
				}
				
				newParam[parts[0]] = val
			} else {
				newParam[parts[0]] = valRaw
			}
		} else {
			if val, ok := postData[key]; ok {
				newParam[key] = val
			} else {
				newParam[key] = 1
			}
		}
	}

	// mulai proses SQL dinamis
	switch method {
	case "insert":
		if enable {
			result := db.Table(tablename).Create(newParam)
			if result.Error != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA CREATE", "TABLE "+tablename)
				return result.Error
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
			postData["lastid"] = fmt.Sprint(lastID)
			helpers.SuccessResponse(c, "DATA SAVED", "")
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
