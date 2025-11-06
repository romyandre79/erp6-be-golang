package generator

import (
	"encoding/json"
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/email"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/i18n"
	"erp6-be-golang/models"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type WorkflowDetailResult struct {
	ComponentDetailID   int    `json:"componentdetailid"`
	ComponentName       string `json:"componentname"`
	DetailType          string `json:"detailtype"`
	Label               string `json:"label"`
	InputType           string `json:"inputtype"`
	InputName           string `json:"inputname"`
	InputDesc           string `json:"inputdesc"`
	DataSourceType      string `json:"datasourcetype"`
	DataSource          string `json:"datasource"`
	DataSourceIDField   string `json:"datasourceidfield"`
	DataSourceNameField string `json:"datasourcenamefield"`
	CompValue           string `json:"compvalue"`
	WfDetailID          int    `json:"wfdetailid"`
}

type WorkflowEngine struct {
	WorkflowId    int
	NodeId        int
	DataInputNode string
	ResultNode    any
}

type Connection struct {
	Node   string `json:"node"`
	Output string `json:"output"`
	Input  string `json:"input"`
}

type IO struct {
	Connections []Connection `json:"connections"`
}

type Component struct {
	WorkflowId int               `json:"workflowid"`
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Class      string            `json:"class"`
	HTML       string            `json:"html"`
	Typenode   bool              `json:"typenode"`
	Inputs     map[string]IO     `json:"inputs"`
	Outputs    map[string]IO     `json:"outputs"`
	PosX       float64           `json:"pos_x"`
	PosY       float64           `json:"pos_y"`
	Data       map[string]string `json:"data"`
	IsRun      bool              `json:"isrun"`
}

type FlowData struct {
	Drawflow struct {
		Home struct {
			Data map[string]Component `json:"data"`
		} `json:"Home"`
	} `json:"drawflow"`
}

func GetWorkflowDetail(db *gorm.DB, componentName string, workflowID int, nodeID int) ([]WorkflowDetailResult, error) {
	var results []WorkflowDetailResult

	var component models.Component
	err := db.
		Preload("Details.WorkflowDetails", "workflowid = ? AND nodeid = ?", workflowID, nodeID).
		Where("componentname = ?", componentName).
		First(&component).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []WorkflowDetailResult{}, nil
		}
		return nil, fmt.Errorf("failed to load component: %w", err)
	}

	for _, d := range component.Details {
		var compValue string
		var wfDetailID int

		if len(d.WorkflowDetails) > 0 {
			compValue = d.WorkflowDetails[0].Componentvalue
			wfDetailID = d.WorkflowDetails[0].Workflowdetailid
		}

		results = append(results, WorkflowDetailResult{
			ComponentDetailID:   int(d.Componentdetailid),
			ComponentName:       component.Componentname,
			DetailType:          d.Detailtype,
			Label:               d.Lable,
			InputType:           d.Inputtype,
			InputName:           d.Inputname,
			InputDesc:           d.Inputdesc,
			DataSourceType:      d.Datasourcetype,
			DataSource:          d.Datasource,
			DataSourceIDField:   d.Datasourceidfield,
			DataSourceNameField: d.Datasourcenamefield,
			CompValue:           compValue,
			WfDetailID:          wfDetailID,
		})
	}

	return results, nil
}

func handleStart(c *fiber.Ctx) error {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: ""})
	c.Locals("wfEngine", wfEngine)
	return nil
}

func handleSendMessage(c *fiber.Ctx, params []WorkflowDetailResult) error {
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	var msg string
	enable := true

	for _, p := range params {
		switch p.InputName {
		case "msgbox":
			msg = strings.TrimSpace(p.CompValue)
		case "enablemsg":
			if strings.EqualFold(p.CompValue, "false") {
				enable = false
			}
		}
	}

	if enable {
		helpers.SuccessResponse(c, "FLOW", msg)
	}
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: msg})
	c.Locals("wfEngine", wfEngine)
	return nil
}

func handleSaveLog(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	logMessage := ""
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	for _, p := range params {
		if p.InputName == "logmessage" {
			logMessage = p.CompValue
		}
	}
	// TODO using log handler
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: logMessage})
	c.Locals("wfEngine", wfEngine)
	return nil
}

func GetSearchText(c *fiber.Ctx, paramTypes []string, param, defVal, dataType string) string {
	s := defVal

	for _, t := range paramTypes {
		switch strings.ToUpper(t) {
		case "POST":
			if val := c.FormValue(param); val != "" {
				s = val
			}
		case "GET":
			if val := c.Query(param); val != "" {
				s = val
			}
		case "Q":
			// Prioritas: GET['q'] > POST['q']
			var val string
			if val = c.Query("q"); val != "" {
				s = val
			} else if val = c.FormValue("q"); val != "" {
				s = val
			}
			fmt.Printf("%s", val)
		}

		// Format tanggal, datetime, time
		if strings.Contains(strings.ToLower(param), "date") && !strings.Contains(strings.ToLower(param), "datetime") {
			if s != "" {
				// ganti "2006-01-02" dengan format ke DB sesuai kebutuhanmu
				t, err := time.Parse("02-01-2006", s) // contoh input: 31-12-2025
				if err == nil {
					s = t.Format("2006-01-02")
				}
			}
		} else if strings.Contains(strings.ToLower(param), "datetime") {
			if s != "" {
				t, err := time.Parse("02-01-2006 15:04:05", s)
				if err == nil {
					s = t.Format("2006-01-02 15:04:05")
				}
			}
		} else if strings.Contains(strings.ToLower(param), "time") {
			if s != "" {
				t, err := time.Parse("15:04:05", s)
				if err == nil {
					s = t.Format("15:04:05")
				}
			}
		}
	}

	// jika tipe data string → ubah ke LIKE pattern: %kata%
	if strings.ToLower(dataType) == "string" {
		s = "%" + strings.ReplaceAll(strings.TrimSpace(s), " ", "%") + "%"
	}

	return s
}

func getUserObjectValues(db *gorm.DB, username string, menuobject string) (string, error) {
	var ids []string
	err := db.Table("groupmenuauth AS a").
		Select("DISTINCT a.menuvalueid").
		Joins("INNER JOIN menuauth b ON b.menuauthid = a.menuauthid").
		Joins("INNER JOIN usergroup c ON c.groupaccessid = a.groupaccessid").
		Joins("INNER JOIN useraccess d ON d.useraccessid = c.useraccessid").
		Where("b.menuobject = ? AND d.username = ?", menuobject, username).
		Pluck("a.menuvalueid", &ids).Error

	if err != nil {
		return "", err
	}

	return strings.Join(ids, ","), nil
}

func getDataByCompany(db *gorm.DB, username string, dataType string) (string, error) {
	cid := ""
	retSql := make(map[string]interface{})
	dataObjectValue, _ := getUserObjectValues(db, username, "company")
	dataSplit := strings.Split(dataObjectValue, ",")
	for _, v := range dataSplit {
		sqlStat := "select a.addressbookid from addressbook a join addresscompany b on b.addressbookid = a.addressbookid where " + dataType + " = 1 and companyid = " + v
		if err := db.Raw(sqlStat).Scan(&retSql).Error; err != nil {
			return "", err
		}
		for _, v := range retSql {
			if cid == "" {
				cid = fmt.Sprintf("%+v", v)
			} else {
				cid += "," + fmt.Sprintf("%+v", v)
			}
		}
	}
	return cid, nil
}

func handleSearch(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	sortStat := ""
	fromStat := ""
	orderStat := ""
	selectStat := ""
	enable := true
	paging := false
	pageStat := 1
	rowsStat := 10
	offsetStat := 0
	leftjoinStat := ""
	whereStat := ""
	total := make(map[string]interface{})
	resultStat := make(map[string]interface{})
	userId := c.Locals("userid")
	userName := c.Locals("username")
	userNameStr, _ := userName.(string)
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)

	for _, p := range params {
		switch p.InputName {
		case "where":
			compValue := strings.TrimSpace(p.CompValue)
			wheres := strings.Fields(compValue)
			for _, data := range wheres {
				dataLower := strings.ToLower(data)
				if dataLower != "and" && dataLower != "or" {
					// Cek apakah ada '@' → berarti fungsi dinamis
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
							// getuserobject atau getuserrecord
							object := strings.Split(funcs[1], ">")
							userObject, _ := getUserObjectValues(db, userNameStr, object[1])
							whereStat += fmt.Sprintf("%s in (%s)", field, userObject)
						case "userid":
							whereStat += fmt.Sprintf("%s = %d", field, userId)
						}
					} else if strings.Contains(data, "=") {
						// Handle ekspresi "="
						datas := strings.SplitN(data, "=", 2)
						left := datas[0]
						right := datas[1]

						if strings.Contains(right, ":") {
							key := strings.ReplaceAll(right, ":", "")
							val := c.Query(key)
							if val == "" {
								val = c.FormValue(key)
							}
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, val)
						} else {
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, right)
						}
					} else {
						// Default → LIKE
						funcs := strings.Split(data, ".")
						whereStat += fmt.Sprintf("(COALESCE(%s,'') LIKE '%s') ", funcs[1], GetSearchText(c, []string{"POST"}, funcs[1], "", "string"))
					}
				} else {
					// and / or
					whereStat += " " + data + " "
				}
			}
		case "paging":
			if p.CompValue == "true" {
				pageStat, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "page", "1", "int"))
				rowsStat, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "rows", "10", "int"))
				offsetStat = (pageStat - 1) * rowsStat
				paging = true
			}
		case "sort":
			if p.CompValue != "" {
				sortStat = strings.Trim(p.CompValue, "")
			}
		case "from":
			if p.CompValue != "" {
				fromStat = strings.Trim(p.CompValue, "")
			}
		case "leftjoin":
			if p.CompValue != "" {
				leftjoinStat = strings.Trim(p.CompValue, "")
			}
		case "select":
			if p.CompValue != "" {
				selectStat = strings.Trim(p.CompValue, "")
			}
		case "enablesearch":
			if p.CompValue == "false" {
				enable = false
			}
		}
	}
	if fromStat != "" {
		if paging {
			sqlStat := "select count(1) as total from " + fromStat
			if leftjoinStat != "" {
				lefts := strings.SplitSeq(leftjoinStat, ",")
				for v := range lefts {
					sqlStat += " left join " + v
				}
			}
			if whereStat != "" {
				sqlStat += " where " + whereStat
			}
			if sortStat != "" {
				sqlStat += " order by " + sortStat + " " + orderStat
			}
			if sqlStat == "" {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "INVALID_QUERY")
			}
			fmt.Printf("sql count %v", sqlStat)

			if enable {
				if err := db.Raw(sqlStat).Scan(&total).Error; err != nil {
					helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
				}
				resultStat = total
			} else {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", sqlStat)
			}
		}

		if selectStat != "" {
			sqlState := "select " + selectStat + " from " + fromStat
			if leftjoinStat != "" {
				lefts := strings.SplitSeq(leftjoinStat, ",")
				for v := range lefts {
					sqlState += " left join " + v
				}
			}
			if whereStat != "" {
				sqlState += " where " + whereStat
			}
			if sortStat != "" {
				sqlState += " order by " + sortStat + " " + orderStat
			}
			if paging {
				sqlState += " limit " + strconv.Itoa(offsetStat) + ", " + strconv.Itoa(rowsStat)
			}
			if sqlState == "" {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "INVALID_QUERY")
			}

			fmt.Printf("sql count %v", sqlState)

			if enable {
				var rows []map[string]interface{}
				if err := db.Raw(sqlState).Scan(&rows).Error; err != nil {
					helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
				}
				resultStat["data"] = rows
				if paging {
					if len(rows) > 0 {
						resultStat["page"] = pageStat
					} else {
						resultStat["page"] = 0
					}
					resultStat["rows"] = rowsStat
				}
				wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
				c.Locals("wfEngine", wfEngine)
				helpers.SuccessResponse(c, "DATA RETRIEVED", resultStat)
			}
		} else {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA RETRIEVED", "EMPTY_QUERY")
		}
	} else {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID DATA RETRIEVED", "EMPTY_QUERY")
	}

	return nil
}

func handleSearchRow(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	sortStat := ""
	fromStat := ""
	orderStat := ""
	selectStat := ""
	enable := true
	leftjoinStat := ""
	whereStat := ""
	resultStat := make(map[string]interface{})
	userId := c.Locals("userid")
	userName := c.Locals("username")
	userNameStr, _ := userName.(string)
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)

	for _, p := range params {
		switch p.InputName {
		case "whererow":
			compValue := strings.TrimSpace(p.CompValue)
			wheres := strings.Fields(compValue)
			for _, data := range wheres {
				dataLower := strings.ToLower(data)
				if dataLower != "and" && dataLower != "or" {
					// Cek apakah ada '@' → berarti fungsi dinamis
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
							// getuserobject atau getuserrecord
							object := strings.Split(funcs[1], ">")
							userObject, _ := getUserObjectValues(db, userNameStr, object[1])
							whereStat += fmt.Sprintf("%s in (%s)", field, userObject)
						case "userid":
							whereStat += fmt.Sprintf("%s = %d", field, userId)
						}
					} else if strings.Contains(data, "=") {
						// Handle ekspresi "="
						datas := strings.SplitN(data, "=", 2)
						left := datas[0]
						right := datas[1]

						if strings.Contains(right, ":") {
							key := strings.ReplaceAll(right, ":", "")
							val := c.Query(key)
							if val == "" {
								val = c.FormValue(key)
							}
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, val)
						} else {
							whereStat += fmt.Sprintf("(COALESCE(%s,'') = '%s') ", left, right)
						}
					} else {
						// Default → LIKE
						funcs := strings.Split(data, ".")
						whereStat += fmt.Sprintf("%s = '%s'", funcs[1], c.FormValue(funcs[1]))
					}
				} else {
					// and / or
					whereStat += " " + data + " "
				}
			}
		case "sortrow":
			if p.CompValue != "" {
				sortStat = strings.Trim(p.CompValue, "")
			}
		case "fromrow":
			if p.CompValue != "" {
				fromStat = strings.Trim(p.CompValue, "")
			}
		case "leftjoinrow":
			if p.CompValue != "" {
				leftjoinStat = strings.Trim(p.CompValue, "")
			}
		case "selectrow":
			if p.CompValue != "" {
				selectStat = strings.Trim(p.CompValue, "")
			}
		case "enablerow":
			if p.CompValue == "false" {
				enable = false
			}
		}
	}

	if selectStat != "" {
		sqlStat := "select " + selectStat + " from " + fromStat
		if leftjoinStat != "" {
			lefts := strings.SplitSeq(leftjoinStat, ",")
			for v := range lefts {
				sqlStat += " left join " + v
			}
		}
		if whereStat != "" {
			sqlStat += " where " + whereStat
		}
		if sortStat != "" {
			sqlStat += " order by " + sortStat + " " + orderStat
		}
		if sqlStat == "" {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "INVALID_QUERY")
		}
		fmt.Printf("sql %s", sqlStat)

		if enable {
			var rows map[string]interface{}
			if err := db.Raw(sqlStat).Scan(&rows).Error; err != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
			}
			if rows != nil {
				resultStat["data"] = rows
				wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
				c.Locals("wfEngine", wfEngine)
				helpers.SuccessResponse(c, "DATA RETRIEVED", resultStat)
			} else {
				helpers.FailResponse(c, 401, "INVALID DATA RETRIEVED", "")
			}
		} else {
			helpers.SuccessResponse(c, "DATA RETRIEVED", sqlStat)
		}
	}

	return nil
}

func handleStoreProcedure(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	var (
		procName string
		paramStr string
		enable   = true
	)

	for _, p := range params {
		switch p.InputName {
		case "paramname":
			paramStr = strings.TrimSpace(p.CompValue)
		case "procedurename":
			procName = strings.TrimSpace(p.CompValue)
		case "enablesp":
			if strings.EqualFold(p.CompValue, "false") {
				enable = false
			}
		}
	}

	if procName == "" {
		return errors.New("INVALID_PROCEDURE_NAME")
	}

	paramList := []string{}
	if paramStr != "" {
		paramList = strings.Split(paramStr, ",")
	}

	var placeholders []string
	for _, p := range paramList {
		placeholders = append(placeholders, ":"+strings.TrimSpace(p))
	}

	sql := fmt.Sprintf("CALL %s(%s)", procName, strings.Join(placeholders, ","))

	if !enable {
		helpers.SuccessResponse(c, "DATA RETRIEVED", sql)
	}

	// Eksekusi stored procedure (dummy)
	if err := db.Exec(sql).Error; err != nil {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
	}
	var wfEngine = c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: paramStr, ResultNode: "OK"})
	c.Locals("wfEngine", wfEngine)
	return nil
}

func handleSendEmail(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	var wfEngine = c.Locals("wfEngine").([]WorkflowEngine)
	emailSender, err := email.NewEmailSender()
	if err != nil {
		helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_EMAIL", err.Error())
	}

	mailTo := ""
	mailHeader := ""
	mailContent := ""
	enable := true

	for _, v := range params {
		switch v.InputName {
		case "mailto":
			mailTo = strings.TrimSpace(v.CompValue)
		case "mailheader":
			mailHeader = strings.TrimSpace(v.CompValue)
		case "mailcontent":
			mailContent = strings.TrimSpace(v.CompValue)
		case "enableemail":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}

	wfEngine = append(wfEngine, WorkflowEngine{
		DataInputNode: "",
		ResultNode:    "OK",
	})
	c.Locals("wfEngine", wfEngine)
	if enable {
		return emailSender.Send(mailTo, mailHeader, mailContent)
	}

	return nil
}

func handleReportServer(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	var (
		vParameter  string
		vReportName string
		vReportType string
	)

	dataPrint := make(map[string]string)

	// Ambil parameter dari workflow
	for _, v := range params {
		switch v.InputName {
		case "reportname":
			vReportName = strings.TrimSpace(v.CompValue)
		case "parameter":
			vParameter = strings.TrimSpace(v.CompValue)
		case "reporttype":
			vReportType = strings.TrimSpace(v.CompValue)
		}
	}

	// Parsing parameter tambahan
	vParams := strings.Split(vParameter, ",")
	for _, v := range vParams {
		dataPrint[v] = GetSearchText(c, []string{"POST"}, v, "", "string")
		dataPrint["title"+v] = v
	}

	lang := c.FormValue("lang")
	userName := c.Locals("username").(string)
	dataPrint["j_username"] = configs.ConfigApps.ReportUser
	dataPrint["j_password"] = configs.ConfigApps.ReportPass
	dataPrint["titlereport"] = i18n.Translate(lang, vReportName, nil)
	dataPrint["titlerecordstatus"] = i18n.Translate(lang, "RECORD_STATUS", nil)
	dataPrint["titlecompany"] = configs.ConfigApps.AppName
	dataPrint["titleuser"] = i18n.Translate(lang, "PRINT_BY", nil) + " " + userName

	timeOut, _ := strconv.Atoi(configs.ConfigApps.ReportTime)
	query := url.Values{}
	for k, v := range dataPrint {
		query.Add(k, v)
	}

	var (
		fullUrl     string
		contentType string
		fileExt     string
	)

	switch strings.ToUpper(vReportType) {
	case "PDF":
		fileExt = "pdf"
		contentType = "application/pdf"
	case "XLS", "XLSX":
		fileExt = "xlsx"
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "CSV":
		fileExt = "csv"
		contentType = "text/csv"
	default:
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_REPORT_TYPE", vReportType)
	}

	fullUrl = fmt.Sprintf("%s/%s.%s?%s", configs.ConfigApps.ReportUrl, vReportName, fileExt, query.Encode())

	data, err := helpers.GetRemoteData(fullUrl, timeOut)
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
	}

	if len(data) == 0 {
		return helpers.FailResponse(c, fiber.StatusNotFound, "EMPTY_REPORT", vReportName)
	}

	// Set header agar file langsung di-download
	c.Set("Cache-Control", "no-cache")
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, vReportName, fileExt))
	c.Set("Content-Length", strconv.Itoa(len(data)))

	return c.Send(data)
}

func cSaveFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}

func handleImportData(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, fileHeader *multipart.FileHeader, search bool) error {
	var (
		column          string
		insertParameter string
		spInsert        string
		updateParameter string
		spUpdate        string
		method          string
		enable          = true
	)

	// 🔹 Ekstrak parameter
	for _, p := range params {
		switch p.InputName {
		case "column":
			column = strings.TrimSpace(p.CompValue)
		case "insertparameter":
			insertParameter = strings.TrimSpace(p.CompValue)
		case "spinsert":
			spInsert = strings.TrimSpace(p.CompValue)
		case "updateparameter":
			updateParameter = strings.TrimSpace(p.CompValue)
		case "importtype":
			method = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "spupdate":
			spUpdate = strings.TrimSpace(p.CompValue)
		case "enableimport":
			if p.CompValue == "false" {
				enable = false
			}
		}
	}

	// 🔹 Simpan file ke "uploads/"
	saveDir := "uploads"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	savePath := filepath.Join(saveDir, filepath.Base(fileHeader.Filename))
	if err := cSaveFile(fileHeader, savePath); err != nil {
		return err
	}

	// 🔹 Baca Excel
	f, err := excelize.OpenFile(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return fmt.Errorf("no data found in Excel")
	}

	// 🔹 Jalankan transaksi (rollback kalau error)
	return db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i < len(rows); i++ {
			colMap := map[string]string{}
			colDefs := strings.Split(column, ",")
			for _, def := range colDefs {
				parts := strings.Split(def, "=")
				if len(parts) != 2 {
					continue
				}
				colName := strings.TrimSpace(parts[0])
				colIndex := strings.TrimSpace(parts[1])
				val, err := f.GetCellValue(sheet, fmt.Sprintf("%s%d", colIndex, i+1))
				if err != nil {
					continue
				}
				colMap[colName] = val
			}

			isEmpty := false
			for _, v := range colMap {
				if strings.TrimSpace(v) == "" {
					isEmpty = true
					break
				}
			}

			var sqlStr string
			var params []string

			if isEmpty {
				params = strings.Split(insertParameter, ",")
				if method == "table" {
					fields, vals := []string{}, []string{}
					for _, p := range params {
						p = strings.TrimSpace(p)
						fields = append(fields, p)
						vals = append(vals, fmt.Sprintf("'%s'", colMap[p]))
					}
					sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", spInsert, strings.Join(fields, ","), strings.Join(vals, ","))
				} else {
					bind := []string{}
					for _, p := range params {
						bind = append(bind, p)
					}
					sqlStr = fmt.Sprintf("CALL %s(%s)", spInsert, strings.Join(bind, ","))
				}
			} else {
				params = strings.Split(updateParameter, ",")
				if method == "table" {
					id := params[0]
					set := []string{}
					for _, p := range params[1:] {
						set = append(set, fmt.Sprintf("%s='%s'", p, colMap[p]))
					}
					sqlStr = fmt.Sprintf("UPDATE %s SET %s WHERE %s='%s'", spUpdate, strings.Join(set, ","), id, colMap[id])
				} else {
					bind := []string{}
					for _, p := range params {
						bind = append(bind, p)
					}
					sqlStr = fmt.Sprintf("CALL %s(%s)", spUpdate, strings.Join(bind, ","))
				}
			}

			if !enable {
				fmt.Println("Preview SQL:", sqlStr)
				continue
			}

			// 🔹 Jalankan dengan GORM
			args := []interface{}{}
			for _, p := range params {
				args = append(args, colMap[p])
			}
			if err := tx.Exec(sqlStr, args...).Error; err != nil {
				log.Printf("Row %d failed: %v", i+1, err)
				return err // rollback
			}
		}
		return nil
	})
}

func handleDecision(c *fiber.Ctx, params []WorkflowDetailResult) (string, bool, error) {
	decisionParamType := ""
	val := ""
	var wfEngine = c.Locals("wfEngine").([]WorkflowEngine)
	enable := true
	contentDecision := ""
	decision := true
	for _, v := range params {
		switch v.InputName {
		case "decisionok":
			contentDecision = v.CompValue
		case "decisionparamtype":
			decisionParamType = strings.ToLower(v.CompValue)
		case "enabledecision":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}

	if strings.Contains(contentDecision, "=") {
		data := strings.SplitN(contentDecision, "=", 2)
		key := data[0]
		expected := data[1]

		switch strings.ToLower(decisionParamType) {
		case "post":
			val = c.FormValue(key)
			if !enable {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "VALUE_DECISION")
			}
			if expected == "empty" {
				decision = val == ""
			} else if val == expected {
				decision = true
			}

		case "get":
			val = c.Query(key)
			if !enable {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "VALUE_DECISION")
			}
			if expected == "empty" {
				decision = val == ""
			} else if val == expected {
				decision = true
			}

		case "node result":
			val = wfEngine[len(wfEngine)-1].ResultNode.(string)
			if !enable {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "VALUE_DECISION")
			}
			if expected == "empty" {
				decision = val == ""
			} else if val == expected {
				decision = true
			}

		default:
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "VALUE_DECISION")
		}
	}

	return contentDecision, decision, nil
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
			newParam[parts[0]] = parts[1]
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

func handleWorkflow(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	wfName := ""
	enable := true
	for _, v := range params {
		switch v.InputName {
		case "workflowname":
			wfName = strings.TrimSpace(v.CompValue)
		//case "workflowparameter":
		//	wfParameter = strings.TrimSpace(v.CompValue)
		case "enableworkflow":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}
	if enable {
		return ExecuteFlow(c, db, wfName, search)
	} else {
		return nil
	}
}

func InternalFlow(c *fiber.Ctx, component Component, workflowId int, nodeId int, db *gorm.DB, search bool) error {
	var flowTerminated = c.Locals("flowTerminated").(bool)
	var components = c.Locals("components").([]Component)
	if flowTerminated {
		return nil
	}

	if component.IsRun {
		return nil
	}

	component.IsRun = true
	fmt.Printf("Running workflowid: %d component: %s (ID: %d)\n", workflowId, component.Name, component.ID)

	workflowDetailResult, err := GetWorkflowDetail(db, component.Name, workflowId, nodeId)
	if err != nil {
		return err
	}

	var Decision bool
	switch component.Name {
	case "Start":
		err = handleStart(c)
	case "SendMessage":
		err = handleSendMessage(c, workflowDetailResult)
	case "SaveLog":
		err = handleSaveLog(c, workflowDetailResult, db, search)
	case "Search":
		err = handleSearch(c, workflowDetailResult, db, search)
	case "SearchRow":
		err = handleSearchRow(c, workflowDetailResult, db, search)
	case "StoreProcedure":
		err = handleStoreProcedure(c, workflowDetailResult, db, search)
	case "SendEmail":
		err = handleSendEmail(c, workflowDetailResult, db, search)
	case "ReportServer":
		err = handleReportServer(c, workflowDetailResult, db, search)
	case "ImportData":
		file, _ := c.FormFile("file-modules")
		err = handleImportData(c, workflowDetailResult, db, file, search)
	case "Decision":
		_, Decision, err = handleDecision(c, workflowDetailResult)
	case "Table":
		err = handleTable(c, workflowDetailResult, db)
	case "Workflow":
		err = handleWorkflow(c, workflowDetailResult, db, search)
	case "End":
		flowTerminated = true
		return nil
	}

	if err != nil {
		return err
	}

	// Jika komponen adalah Decision
	if component.Name == "Decision" {
		var outputs IO
		if Decision {
			outputs = component.Outputs["output_1"]
		} else {
			outputs = component.Outputs["output_2"]
		}

		for _, conn := range outputs.Connections {
			nextNodeId, _ := strconv.Atoi(conn.Node)
			for _, nextComp := range components {
				if nextComp.ID == nextNodeId && nextComp.WorkflowId == workflowId {
					return InternalFlow(c, nextComp, workflowId, nextNodeId, db, search)
				}
			}
		}
		return nil
	} else {

		// Untuk node biasa, lanjut ke semua output berurutan
		for _, output := range component.Outputs {
			for _, conn := range output.Connections {
				nextNodeId, _ := strconv.Atoi(conn.Node)
				for _, nextComp := range components {
					if nextComp.ID == nextNodeId && nextComp.WorkflowId == workflowId {
						err := InternalFlow(c, nextComp, workflowId, nextNodeId, db, search)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func ExecuteFlow(c *fiber.Ctx, db *gorm.DB, flowName string, search bool) error {
	var components = []Component{}
	var wfEngine = []WorkflowEngine{}
	var flowTerminated = false

	// === 1️⃣ Ambil parameter workflow ===
	var params []models.Workflowparameter
	if err := db.
		Table("workflowparameter a").
		Select("a.wfparameterid, b.wfname, a.parametername").
		Joins("INNER JOIN workflow b ON b.workflowid = a.workflowid").
		Where("b.wfname = ?", flowName).
		Scan(&params).Error; err != nil {
		return err
	}

	// === 2️⃣ Inisialisasi parameter default ===
	postData := make(map[string]interface{})
	for _, p := range params {
		postData[p.Parametername] = nil
	}

	// === 3️⃣ Ambil definisi workflow utama ===
	var wf models.Workflow
	if err := db.
		Where("wfname = ?", flowName).
		First(&wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("workflow %s does not exist", flowName)
		}
		return errors.New("INVALID_WORKFLOW")
	}

	// === 4️⃣ Decode JSON flow ke struktur data ===
	var flow FlowData
	if err := json.Unmarshal([]byte(wf.Flow), &flow); err != nil {
		return fmt.Errorf("failed to decode flow JSON: %w", err)
	}

	// === 5️⃣ Konversi map → slice agar bisa diurutkan ===
	for _, comp := range flow.Drawflow.Home.Data {
		comp.WorkflowId = wf.Workflowid
		components = append(components, comp)
	}

	// === 6️⃣ Bangun graph koneksi antar node ===
	graph := make(map[int][]int)
	inDegree := make(map[int]int)
	for _, comp := range components {
		for _, out := range comp.Outputs {
			for _, conn := range out.Connections {
				src := comp.ID
				dst, _ := strconv.Atoi(conn.Node)
				graph[src] = append(graph[src], dst)
				inDegree[dst]++
				if _, ok := inDegree[src]; !ok {
					inDegree[src] = 0
				}
			}
		}
	}

	// === 7️⃣ Temukan node awal (Start node: inDegree == 0) ===
	queue := []int{}
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	// === 8️⃣ Urutkan node secara topologis ===
	ordered := []Component{}
	visited := make(map[int]bool)

	for len(queue) > 0 {
		currID := queue[0]
		queue = queue[1:]

		if visited[currID] {
			continue
		}
		visited[currID] = true

		for _, comp := range components {
			if comp.ID == currID {
				comp.WorkflowId = wf.Workflowid
				ordered = append(ordered, comp)
				break
			}
		}

		for _, next := range graph[currID] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	// === 9️⃣ Jika ada node tidak terhubung, tambahkan ke akhir ===
	if len(ordered) < len(components) {
		for _, comp := range components {
			if !visited[comp.ID] {
				ordered = append(ordered, comp)
			}
		}
	}

	// === 🔟 Debug: tampilkan urutan flow ===
	//out, _ := json.MarshalIndent(ordered, "", "  ")
	//fmt.Println("Ordered flow:", string(out))

	// === 11️⃣ Siapkan transaksi database ===
	var tx *gorm.DB
	if !search {
		tx = db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				fmt.Printf("Transaction panic: %v\n", r)
			}
		}()
	} else {
		tx = db
	}

	c.Locals("components", ordered)
	c.Locals("wfEngine", wfEngine)
	c.Locals("flowTerminated", flowTerminated)

	// === 12️⃣ Eksekusi tiap komponen sesuai urutan ===
	//for _, comp := range ordered {
	if err := InternalFlow(c, ordered[0], int(wf.Workflowid), ordered[0].ID, tx, search); err != nil {
		if !search {
			tx.Rollback()
		}
		return fmt.Errorf("internal flow failed on component %d (%s): %w", ordered[0].ID, ordered[0].Name, err)
	}
	//}

	// === 13️⃣ Commit transaksi ===
	if !search {
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
	}

	return nil
}

// ============================
// Helper
// ============================
func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
