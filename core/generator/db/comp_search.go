package generator

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("search", func(ctx *WorkflowContext) error {
		return handleSearch(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
	RegisterComponent("searchrow", func(ctx *WorkflowContext) error {
		return handleSearchRow(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
	RegisterComponent("searchsingle", func(ctx *WorkflowContext) error {
		return handleSearchSingle(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
}

// SearchParams holds the parsed parameters for a search operation
type SearchParams struct {
	Sort     string
	From     string
	Order    string
	Select   string
	LeftJoin string
	Where    string
	GroupBy  string
	Enable   bool
	Paging   bool
	Page     int
	Rows     int
	Offset   int
}

func parseWhereClause(c *fiber.Ctx, db *gorm.DB, compValue string, userNameStr string, userId interface{}, isRow bool) string {
	driver := GetDatabaseDriver(db)
	whereStat := ""
	wheres := strings.Fields(compValue)
	for i := 0; i < len(wheres); i++ {
		data := wheres[i]

		// Resolve any $parameters within the token (e.g., month($postdate) -> month('2024-07-04'))
		// This handles SQL functions containing parameters
		if strings.Contains(data, "$") {
			data = ResolveParam(c, data)
		}

		dataLower := strings.ToLower(data)

		// Check for lookahead "IN" or "NOT IN" operator to prevent auto-LIKE on dot-notation fields
		// AND to handle "IN *" or "NOT IN *" (wildcard checks)
		isNextIn := false

		// If this token is a SQL operator, pass it through as-is
		if data == "=" || data == "!=" || data == ">" || data == "<" || data == ">=" || data == "<=" || data == "<>" {
			whereStat += " " + data + " "
			continue
		}

		// If this token is a SQL function (contains parentheses), pass it through as-is
		// This prevents month(2023-07-04) from being treated as a field name
		if strings.Contains(data, "(") && strings.Contains(data, ")") {
			whereStat += " " + data + " "
			continue
		}

		// Lookahead for Field LIKE Value
		if i+2 < len(wheres) {
			if strings.ToLower(wheres[i+1]) == "like" {
				// field like $param
				field := wheres[i]
				valParam := wheres[i+2]

				// Resolve value
				val := ResolveParam(c, valParam)
				val = strings.Trim(val, "'") // Strip quotes if present

				if val == "null" {
					val = ""
				}

				// If value is still a variable (unresolved), treat as empty
				if strings.HasPrefix(val, "$") {
					val = ""
				}

				// Wrap in % for like
				val = "%" + val + "%"

				// Construct SQL
				whereStat += fmt.Sprintf("(COALESCE(%s, '') LIKE '%s') ", field, val)

				i += 2
				continue
			}
		}

		// Lookahead for Field IN Value
		if i+2 < len(wheres) {
			if strings.ToLower(wheres[i+1]) == "in" {
				isNextIn = true

				// Check if Value is wildcard '*'
				val := ResolveParam(c, wheres[i+2])
				cleanVal := strings.Trim(val, "()")

				if cleanVal == "*" {
					// Field IN * -> Treat as ALL (True)
					whereStat += " 1=1 "
					i += 2 // Skip Field, IN, Value
					continue
				}
			}
		}

		// Lookahead for Field NOT IN Value
		if i+3 < len(wheres) {
			if strings.ToLower(wheres[i+1]) == "not" && strings.ToLower(wheres[i+2]) == "in" {
				isNextIn = true

				// Check if Value is wildcard '*'
				val := ResolveParam(c, wheres[i+3])
				cleanVal := strings.Trim(val, "()")

				if cleanVal == "*" {
					// Field NOT IN * -> Treat as NONE (False)
					whereStat += " 0=1 "
					i += 3 // Skip Field, NOT, IN, Value
					continue
				}
			}
		}

		if dataLower != "and" && dataLower != "or" {
			// Check if there's '@' → dynamic function
			if strings.Contains(data, "@") {
				funcs := strings.Split(data, "@")
				field := strings.ReplaceAll(funcs[0], "@", "")
				field = strings.ReplaceAll(field, "=", "")

				switch strings.ToLower(funcs[1]) {
				case "getemployeebycompany":
					listByCompany, _ := getDataByCompany(db, userNameStr, "isemployee")
					whereStat += fmt.Sprintf("%s in (%s)", field, listByCompany)
				case "getcustomerbycompany":
					listByCompany, _ := getDataByCompany(db, userNameStr, "iscustomer")
					whereStat += fmt.Sprintf("%s in (%s)", field, listByCompany)
				case "getvendorbycompany":
					listByCompany, _ := getDataByCompany(db, userNameStr, "isvendor")
					whereStat += fmt.Sprintf("%s in (%s)", field, listByCompany)
				case "getuserrecord":
					// getuserobject or getuserrecord
					object := strings.Split(funcs[1], ">")
					userObject, _ := getUserObjectValues(db, userNameStr, object[1])
					whereStat += fmt.Sprintf("%s in (%s)", field, userObject)
				case "userid":
					whereStat += fmt.Sprintf("%s = %v", field, userId)
				}
			} else if strings.Contains(data, "!=") {
				datas := strings.SplitN(data, "!=", 2)
				left := datas[0]
				right := datas[1]

				if strings.Contains(right, "$") {
					key := strings.ReplaceAll(right, "$", "")
					val := c.Query(key)
					if val == "" {
						val = c.FormValue(key)
					}
					switch driver {
					case "mysql", "mariadb":
						whereStat += fmt.Sprintf("(COALESCE(%s,'') <> '%s') ", left, val)
					default:
						whereStat += fmt.Sprintf("(COALESCE(%s,'') != '%s') ", left, val)
					}
				} else {
					if strings.Contains(data, "empty") {
						whereStat += fmt.Sprintf("%s is not null ", left)
					} else {
						switch driver {
						case "mysql", "mariadb":
							whereStat += fmt.Sprintf("(COALESCE(%s,'') <> '%s') ", left, right)
						default:
							whereStat += fmt.Sprintf("(COALESCE(%s,'') != '%s') ", left, right)
						}
					}
				}
			} else if strings.Contains(data, "=") {
				datas := strings.SplitN(data, "=", 2)
				left := datas[0]
				right := datas[1]

				if strings.HasPrefix(strings.ToLower(right), "between:") {
					// Format: field=between:start:end
					parts := strings.Split(right, ":")
					if len(parts) == 3 {
						start := parts[1]
						end := parts[2]

						// Check if values are dynamic (from query/form)
						if strings.HasPrefix(start, "$") {
							key := start[1:]
							val := c.Query(key)
							if val == "" {
								val = c.FormValue(key)
							}
							start = val
						}
						if strings.HasPrefix(end, "$") {
							key := end[1:]
							val := c.Query(key)
							if val == "" {
								val = c.FormValue(key)
							}
							end = val
						}

						whereStat += fmt.Sprintf("(%s BETWEEN '%s' AND '%s') ", left, start, end)
					}
				} else if strings.Contains(right, "$") {
					key := strings.ReplaceAll(right, "$", "")
					val := c.Query(key)
					if val == "" {
						val = c.FormValue(key)
					}
					whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, val)
				} else {
					if strings.Contains(data, "empty") {
						whereStat += fmt.Sprintf("%s is null", left)
					} else if strings.Contains(data, "exist") {
						whereStat += fmt.Sprintf("exist (%s)", left)
					} else {
						if strings.Contains(right, "$") {
							key := strings.ReplaceAll(right, "$", "")
							val := c.Query(key)
							if val == "" {
								val = c.FormValue(key)
							}
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, val)
						} else {
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, right)
						}
					}
				}
			} else {
				// Default → LIKE
				if strings.Contains(data, "=") {
					funcs := strings.Split(data, "=")
					val := GetSearchText(c, []string{"POST"}, funcs[1], "", "string")
					whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", data, val)
				} else if strings.Contains(data, ".") && !isNextIn { // Skip auto-LIKE if followed by IN
					funcs := strings.Split(data, ".")
					val := GetSearchText(c, []string{"POST"}, funcs[1], "", "string")
					cleanVal := strings.ReplaceAll(val, "%", "")
					if strings.HasSuffix(strings.ToLower(funcs[1]), "id") {
						// Strip wildcards for ID exact match
						if cleanVal == "" {
							if isRow {
								whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", data, cleanVal)
							} else {
								whereStat += fmt.Sprintf("(COALESCE(%s,'') LIKE '%s') ", data, val)
							}
						} else {
							whereStat += fmt.Sprintf("(%s = '%s') ", data, cleanVal)
						}
					} else {
						if isRow {
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", data, cleanVal)
						} else {
							whereStat += fmt.Sprintf("(COALESCE(%s,'') LIKE '%s') ", data, val)
						}
					}
				} else {
					whereStat += " " + ResolveParam(c, data) + " "
				}
			}
		} else {
			// and / or
			whereStat += " " + ResolveParam(c, data) + " "
		}
	}
	return whereStat
}

func parseSelectClause(selectStr string) string {
	parts := strings.Split(selectStr, ",")
	var newParts []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, ":") {
			subParts := strings.Split(part, ":")
			if len(subParts) == 2 {
				field := subParts[0]
				op := strings.ToLower(subParts[1])
				switch op {
				case "count":
					if field == "*" || field == "1" {
						newParts = append(newParts, fmt.Sprintf("COUNT(%s)", field))
					} else {
						newParts = append(newParts, fmt.Sprintf("COUNT(%s)", field))
					}
				case "sum":
					newParts = append(newParts, fmt.Sprintf("SUM(%s)", field))
				case "avg":
					newParts = append(newParts, fmt.Sprintf("AVG(%s)", field))
				case "min":
					newParts = append(newParts, fmt.Sprintf("MIN(%s)", field))
				case "max":
					newParts = append(newParts, fmt.Sprintf("MAX(%s)", field))
				default:
					newParts = append(newParts, part)
				}
			} else if len(subParts) == 3 {
				// field:op:alias
				field := subParts[0]
				op := strings.ToLower(subParts[1])
				alias := subParts[2]
				switch op {
				case "count":
					if field == "*" || field == "1" {
						newParts = append(newParts, fmt.Sprintf("COUNT(%s) as %s", field, alias))
					} else {
						newParts = append(newParts, fmt.Sprintf("COUNT(%s) as %s", field, alias))
					}
				case "sum":
					newParts = append(newParts, fmt.Sprintf("SUM(%s) as %s", field, alias))
				case "avg":
					newParts = append(newParts, fmt.Sprintf("AVG(%s) as %s", field, alias))
				case "min":
					newParts = append(newParts, fmt.Sprintf("MIN(%s) as %s", field, alias))
				case "max":
					newParts = append(newParts, fmt.Sprintf("MAX(%s) as %s", field, alias))
				default:
					newParts = append(newParts, part)
				}
			} else {
				newParts = append(newParts, part)
			}
		} else {
			newParts = append(newParts, part)
		}
	}
	return strings.Join(newParts, ", ")
}

func parseSearchParams(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, suffix string, isRow bool) SearchParams {
	sp := SearchParams{
		Enable: true,
		Page:   1,
		Rows:   10,
	}
	userId := c.Locals("userid")
	userName := c.Locals("username")
	userNameStr, _ := userName.(string)

	for _, p := range params {
		name := strings.ToLower(p.InputName)
		val := strings.TrimSpace(p.CompValue)

		if name == "where"+suffix {
			sp.Where = parseWhereClause(c, db, val, userNameStr, userId, isRow)
		} else if name == "paging" && suffix == "" { // Paging usually only for main search
			if val == "true" {
				sp.Page, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "page", "1", "int"))
				sp.Rows, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "rows", "10", "int"))
				sp.Offset = (sp.Page - 1) * sp.Rows
				if sp.Offset < 0 {
					sp.Offset = 0
				}
				sp.Paging = true
			}
		} else if name == "sort"+suffix {
			sp.Sort = val
		} else if name == "from"+suffix {
			sp.From = val
		} else if name == "leftjoin"+suffix {
			sp.LeftJoin = val
		} else if name == "select"+suffix {
			sp.Select = parseSelectClause(val)
		} else if name == "groupby"+suffix {
			sp.GroupBy = val
		} else if name == "enable"+suffix {
			if val == "false" {
				sp.Enable = false
			}
		}
	}

	// Dynamic sorting from request (overrides workflow config)
	if sortBy := c.FormValue("sortby"); sortBy != "" {
		sp.Sort = sortBy
	} else if sortBy := c.Query("sortby"); sortBy != "" {
		sp.Sort = sortBy
	}

	if sortDir := c.FormValue("sortdir"); sortDir != "" {
		// Validate sortdir to prevent SQL injection
		sortDirLower := strings.ToLower(sortDir)
		if sortDirLower == "asc" || sortDirLower == "desc" {
			sp.Order = sortDirLower
		}
	} else if sortDir := c.Query("sortdir"); sortDir != "" {
		sortDirLower := strings.ToLower(sortDir)
		if sortDirLower == "asc" || sortDirLower == "desc" {
			sp.Order = sortDirLower
		}
	}

	return sp
}

func handleGenericSearch(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, suffix string, isSingle bool, isRow bool) error {
	sp := parseSearchParams(c, params, db, suffix, isRow)
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	resultStat := make(map[string]interface{})

	if sp.From == "" {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA RETRIEVED", "EMPTY_QUERY")
		return nil
	}

	// Count query for paging
	if sp.Paging && !isSingle && !isRow {
		var sqlStat string

		// When GROUP BY is used, count the number of groups
		if sp.GroupBy != "" {
			innerQuery := "select 1 from " + sp.From
			if sp.LeftJoin != "" {
				lefts := strings.Split(sp.LeftJoin, ",")
				for _, v := range lefts {
					innerQuery += " left join " + v
				}
			}
			if sp.Where != "" {
				innerQuery += " where " + sp.Where
			}
			innerQuery += " group by " + sp.GroupBy
			sqlStat = fmt.Sprintf("select count(1) as total from (%s) as grouped_results", innerQuery)
		} else {
			sqlStat = "select count(1) as total from " + sp.From
			if sp.LeftJoin != "" {
				lefts := strings.Split(sp.LeftJoin, ",")
				for _, v := range lefts {
					sqlStat += " left join " + v
				}
			}
			if sp.Where != "" {
				sqlStat += " where " + sp.Where
			}
		}

		if sp.Enable {
			var total map[string]interface{}
			if err := db.Raw(sqlStat).Scan(&total).Error; err != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
				return nil
			}
			resultStat = total
		} else {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", sqlStat)
			return nil
		}
	}

	if sp.Select != "" {
		sqlState := "select " + sp.Select + " from " + sp.From
		if sp.LeftJoin != "" {
			lefts := strings.Split(sp.LeftJoin, ",")
			for _, v := range lefts {
				sqlState += " left join " + v
			}
		}
		if sp.Where != "" {
			sqlState += " where " + sp.Where
		}
		if sp.GroupBy != "" {
			sqlState += " group by " + sp.GroupBy
		}
		if sp.Sort != "" {
			sqlState += " order by " + sp.Sort + " " + sp.Order
		}
		if sp.Paging && !isSingle && !isRow {
			driver := GetDatabaseDriver(db)
			switch driver {
			case "sqlserver", "oracle":
				// SQL Server and Oracle 12c+ support OFFSET/FETCH
				if sp.Sort == "" {
					sqlState += " order by (SELECT NULL)"
				}
				sqlState += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", sp.Offset, sp.Rows)
			default:
				// Work for MySQL, PostgreSQL, SQLite
				// LIMIT rows OFFSET offset
				sqlState += fmt.Sprintf(" LIMIT %d OFFSET %d", sp.Rows, sp.Offset)
			}
		}

		log.Info(sqlState)

		if sp.Enable {
			if isSingle {
				var singleResult string
				if err := db.Raw(sqlState).Scan(&singleResult).Error; err != nil {
					helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
					return nil
				}
				resultStat["data"] = singleResult
			} else if isRow {
				var rowResult map[string]interface{}
				if err := db.Raw(sqlState).Scan(&rowResult).Error; err != nil {
					helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
					return nil
				}
				if rowResult != nil {
					resultStat["data"] = rowResult
				} else {
					helpers.SuccessResponse(c, "INVALID DATA RETRIEVED", "")
					return nil
				}
			} else {
				var rows []map[string]interface{}
				if err := db.Raw(sqlState).Scan(&rows).Error; err != nil {
					helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
					return nil
				}
				resultStat["data"] = rows
				if sp.Paging {
					if len(rows) > 0 {
						if sp.Page == 0 {
							resultStat["page"] = 1
						} else {
							resultStat["page"] = sp.Page
						}
					} else {
						resultStat["page"] = 0
					}
					resultStat["rows"] = sp.Rows
				}
			}

			wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
			c.Locals("wfEngine", wfEngine)
			helpers.SuccessResponse(c, "DATA RETRIEVED", resultStat)
		} else {
			helpers.SuccessResponse(c, "DATA RETRIEVED", sqlState)
		}
	} else {
		helpers.SuccessResponse(c, "INVALID DATA RETRIEVED", "EMPTY_QUERY")
	}

	return nil
}

func handleSearch(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	return handleGenericSearch(c, params, db, "", false, false)
}

func handleSearchRow(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	return handleGenericSearch(c, params, db, "row", false, true)
}

func handleSearchSingle(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	return handleGenericSearch(c, params, db, "single", true, false)
}
