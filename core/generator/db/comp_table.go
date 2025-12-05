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
		return handleTable(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

func handleTable(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
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

	listOldParam := strings.Split(param, ",")
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
			val := c.Query(parts[1])
			if val == "" {
				val = c.FormValue(parts[1])
			}
			newParam[parts[0]] = val
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
			db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
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
