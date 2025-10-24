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

type Component struct {
	ID      string
	Name    string
	Outputs map[string]interface{}
}

type WorkflowEngine struct {
	WorkflowId    int
	NodeId        int
	DataInputNode string
	ResultNode    any
}

var postData = make(map[string]interface{})
var components = []Component{}
var wfEngine = []WorkflowEngine{}

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

func handleStart(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	fmt.Print("Start")

	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: ""})
	return nil
}

func handleSendMessage(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	fmt.Print("Send Message")
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
	return nil
}

func handleSaveLog(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	fmt.Print("Save Log")
	logMessage := ""
	for _, p := range params {
		if p.InputName == "logmessage" {
			logMessage = p.CompValue
		}
	}
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: logMessage})
	fmt.Printf("[LOG] %s\n", logMessage)
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
			if val := c.Query("q"); val != "" {
				s = val
			} else if val := c.FormValue("q"); val != "" {
				s = val
			}
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
	fmt.Print("Search")
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

	for _, p := range params {
		fmt.Printf("Params %+v\n", p)
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
						userName := c.Locals("username")
						userNameStr, _ := userName.(string)

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
						}
					} else if strings.Contains(data, "=") {
						// Handle ekspresi "="
						datas := strings.SplitN(data, "=", 2)
						fmt.Printf("datas %+v", datas)
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
						whereStat += fmt.Sprintf("(COALESCE(%s,'') LIKE '%s') ", data, GetSearchText(c, []string{"Q"}, data, "", "string"))
					}
				} else {
					// and / or
					whereStat += " " + data + " "
				}
			}
		case "paging":
			pageStat, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "page", "1", "int"))
			rowsStat, _ = strconv.Atoi(GetSearchText(c, []string{"POST", "GET"}, "rows", "10", "int"))
			offsetStat = (pageStat - 1) * rowsStat
			paging = true
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
	sqlStat := "select count(1) as total from " + fromStat
	if leftjoinStat != "" {
		lefts := strings.Split(leftjoinStat, ",")
		for _, v := range lefts {
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
	if enable {
		if err := db.Raw(sqlStat).Scan(&total).Error; err != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
		}
		resultStat = total
	} else {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", sqlStat)
	}

	sqlStat = "select " + selectStat + " from " + fromStat
	if leftjoinStat != "" {
		lefts := strings.Split(leftjoinStat, ",")
		for _, v := range lefts {
			sqlStat += " left join " + v
		}
	}
	if whereStat != "" {
		sqlStat += " where " + whereStat
	}
	if sortStat != "" {
		sqlStat += " order by " + sortStat + " " + orderStat
	}
	if paging {
		sqlStat += " limit " + strconv.Itoa(offsetStat) + ", " + strconv.Itoa(rowsStat)
	}
	if sqlStat == "" {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", "INVALID_QUERY")
	}

	if enable {
		var rows []map[string]interface{}
		if err := db.Raw(sqlStat).Scan(&rows).Error; err != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
		}
		resultStat["data"] = rows
		resultStat["page"] = GetSearchText(c, []string{"POST", "GET"}, "page", "1", "int")
		resultStat["rows"] = GetSearchText(c, []string{"POST", "GET"}, "rows", "10", "int")
		wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultStat})
		helpers.SuccessResponse(c, "DATA_RETRIEVED", resultStat)
	} else {
		helpers.SuccessResponse(c, "DATA_RETRIEVED", sqlStat)
	}

	return nil
}

func handleStoreProcedure(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	fmt.Print("Store Procedure")
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
		helpers.SuccessResponse(c, "DATA_RETRIEVED", sql)
	}

	// Eksekusi stored procedure (dummy)
	if err := db.Exec(sql).Error; err != nil {
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
	}

	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: paramStr, ResultNode: "OK"})
	return nil
}

func handleSendEmail(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
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
	if enable {
		return emailSender.Send(mailTo, mailHeader, mailContent)
	}

	return nil
}

func ReportPrint(c *fiber.Ctx, printType string, reportName string, dataPrint map[string]string) error {
	lang := c.FormValue("lang")
	userName := c.Locals("username").(string)
	dataPrint["j_username"] = configs.ConfigApps.ReportUser
	dataPrint["j_password"] = configs.ConfigApps.ReportPass
	dataPrint["titlereport"] = i18n.Translate(lang, reportName, nil)
	dataPrint["titlerecordstatus"] = i18n.Translate(lang, "RECORD_STATUS", nil)
	dataPrint["titlecompany"] = configs.ConfigApps.AppName
	dataPrint["titleuser"] = i18n.Translate(lang, "PRINT_BY", nil) + " " + userName
	timeOut, _ := strconv.Atoi(configs.ConfigApps.ReportTime)

	data := []byte{}
	query := url.Values{}
	for k, v := range dataPrint {
		query.Add(k, v)
	}

	switch printType {
	case "PDF":
		fullUrl := fmt.Sprintf("%s/%s.pdf?%s", configs.ConfigApps.ReportUrl, reportName, query.Encode())
		data, err := helpers.GetRemoteData(fullUrl, timeOut)
		if err != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
		}
		if strings.HasPrefix(string(data), "%PDF") {
			c.Set("Cache-Control", "public")
			c.Set("Content-Type", "application/pdf")
			c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.pdf"`, reportName))
			c.Set("Content-Length", strconv.Itoa(len(data)))
			return c.Send(data)
		}
	case "XLS":
		fullUrl := fmt.Sprintf("%s/%s.xlsx?%s", configs.ConfigApps.ReportUrl, reportName, query.Encode())
		data, err := helpers.GetRemoteData(fullUrl, timeOut)
		if err != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
		}
		c.Set("Cache-Control", "public")
		c.Set("Content-Type", "application/vnd.ms-excel")
		c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.xlsx"`, reportName))
		c.Set("Content-Length", strconv.Itoa(len(data)))
		return c.Send(data)
	case "CSV":
		fullUrl := fmt.Sprintf("%s/%s.csv?%s", configs.ConfigApps.ReportUrl, reportName, query.Encode())
		data, err := helpers.GetRemoteData(fullUrl, timeOut)
		if err != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_FLOW", err.Error())
		}
		c.Set("Cache-Control", "public")
		c.Set("Content-Type", "application/vnd.ms-word")
		c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.csv"`, reportName))
		c.Set("Content-Length", strconv.Itoa(len(data)))
		return c.Send(data)
	}
	return c.SendString(string(data))
}

func handleReportServer(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	vParameter := ""
	vReportName := ""
	vReportType := ""
	enable := true
	dataPrint := make(map[string]string)
	for _, v := range params {
		switch v.InputName {
		case "reportname":
			vReportName = strings.TrimSpace(v.CompValue)
		case "parameter":
			vParameter = strings.TrimSpace(v.CompValue)
		case "reporttype":
			vReportType = strings.TrimSpace(v.CompValue)
		case "enablereport":
			if v.CompValue == "false" {
				enable = false
			}
		}
	}
	vParams := strings.Split(vParameter, ",")
	lang := c.FormValue("lang")
	for _, v := range vParams {
		dataPrint[v] = GetSearchText(c, []string{"GET"}, v, "", "string")
		dataPrint["title"+v] = i18n.Translate(lang, v, nil)
	}
	if enable {
		return ReportPrint(c, vReportType, vReportName, dataPrint)
	}
	return nil
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

func handleDecision(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) (string, bool, error) {
	decisionParamType := ""
	val := ""
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

func handleTable(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
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
	newParam := map[string]string{}
	for _, key := range listOldParam {
		if strings.Contains(key, "=") {
			parts := strings.SplitN(key, "=", 2)
			newParam[parts[0]] = parts[1]
		} else {
			if val, ok := postData[key]; ok {
				newParam[key] = val
			} else {
				newParam[key] = "1"
			}
		}
	}

	// mulai proses SQL dinamis
	switch method {
	case "insert":
		if enable {
			result := db.Table(tablename).Create(newParam)
			if result.Error != nil {
				return result.Error
			}

			var lastID int64
			db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
			postData["lastid"] = fmt.Sprint(lastID)
		}
		helpers.SuccessResponse(c, "DATA_SAVED", "TABLE "+tablename)

	case "update":
		idField := listOldParam[0]
		idValue := postData[idField]
		delete(newParam, idField)

		if enable {
			result := db.Table(tablename).Where(fmt.Sprintf("%s = ?", idField), idValue).Updates(newParam)
			if result.Error != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_DATA_UPDATE", "TABLE "+tablename)
			}
		}
		helpers.SuccessResponse(c, "DATA_SAVED", "TABLE "+tablename)

	case "purge":
		idField := listOldParam[0]
		idValue := postData["id"]

		if enable {
			result := db.Table(tablename).Where(fmt.Sprintf("%s = ?", idField), idValue).Delete(nil)
			if result.Error != nil {
				helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_DATA_PURGE", "TABLE "+tablename)
			}
		}
		helpers.SuccessResponse(c, "DATA_SAVED", "TABLE "+tablename)

	default:
		helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_DATA_UPDATE", "TABLE "+tablename)
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
	fmt.Printf("component %s", component.Name)
	Decision := false

	workflowDetailResult, err := GetWorkflowDetail(db, component.Name, workflowId, nodeId)
	if err != nil {
		return err
	}

	err = nil

	switch component.Name {
	case "Start":
		err = handleStart(c, workflowDetailResult, db, search)
	case "SendMessage":
		err = handleSendMessage(c, workflowDetailResult, db, search)
	case "SaveLog":
		err = handleSaveLog(c, workflowDetailResult, db, search)
	case "Search":
		err = handleSearch(c, workflowDetailResult, db, search)
	case "StoreProcedure":
		err = handleStoreProcedure(c, workflowDetailResult, db, search)
	case "SendEmail":
		err = handleSendEmail(c, workflowDetailResult, db, search)
	case "ReportServer":
		err = handleReportServer(c, workflowDetailResult, db, search)
	case "ImportData":
		file, err := c.FormFile("file-modules")
		if err != nil {
			return err
		}
		err = handleImportData(c, workflowDetailResult, db, file, search)
		if err != nil {
			return err
		}
	case "Decision":
		_, Decision, err = handleDecision(c, workflowDetailResult, db, search)
	case "Table":
		err = handleTable(c, workflowDetailResult, db, search)
	case "Workflow":
		err = handleWorkflow(c, workflowDetailResult, db, search)
	}

	if component.Name == "Search" && search {
		return err
	}

	if component.Name == "Decision" {
		if Decision {
			outputs, ok := component.Outputs["output_1"].([]interface{})
			if !ok {
				return errors.New("")
			}

			for _, row2 := range outputs {
				rowSlice, ok := row2.([]interface{})
				if !ok || len(rowSlice) == 0 {
					continue
				}

				// ambil node dari row2[0]["node"]
				rowMap, ok := rowSlice[0].(map[string]interface{})
				if !ok {
					continue
				}
				node, _ := rowMap["node"].(string)
				nodeId, _ := strconv.Atoi(node)

				for _, comp := range components {
					if comp.ID == node {
						InternalFlow(c, comp, workflowId, nodeId, db, search)
					}
				}
			}
		} else {
			outputs, ok := component.Outputs["output_2"].([]interface{})
			if !ok {

			}

			for _, row2 := range outputs {
				rowSlice, ok := row2.([]interface{})
				if !ok || len(rowSlice) == 0 {
					continue
				}

				// ambil node dari row2[0]["node"]
				rowMap, ok := rowSlice[0].(map[string]interface{})
				if !ok {
					continue
				}
				node, _ := rowMap["node"].(string)
				nodeId, _ := strconv.Atoi(node)

				for _, comp := range components {
					if comp.ID == node {
						InternalFlow(c, comp, workflowId, nodeId, db, search)
					}
				}
			}
		}
	} else {
		for _, row1Any := range component.Outputs {
			row1, ok := row1Any.(map[string]interface{})
			if !ok {
				break
			}

			connections, ok := row1["connections"].([]interface{})
			if !ok {
				break
			}

			for _, row2Any := range connections {
				row2, ok := row2Any.(map[string]interface{})
				if !ok {
					break
				}

				node, _ := row2["node"].(string)
				nodeIdr, _ := strconv.Atoi(node)

				for _, comp := range components {
					if comp.ID == node {
						return InternalFlow(c, comp, workflowId, nodeIdr, db, search)
					}
				}
			}
		}
	}

	return nil
}

func ExecuteFlow(c *fiber.Ctx, db *gorm.DB, flowName string, search bool) error {
	var params []models.Workflowparameter
	if err := db.
		Table("workflowparameter a").
		Select("a.wfparameterid, b.wfname, a.parametername").
		Joins("inner join workflow b on b.workflowid = a.workflowid").
		Where("b.wfname = ?", flowName).
		Scan(&params).Error; err != nil {
		fmt.Errorf("failed to get workflow parameters: %s", err.Error())
		return err
	}

	for _, p := range params {
		fmt.Printf("params %v", p)
		postData[p.Parametername] = nil
	}

	var wf models.Workflow
	if err := db.
		Where("wfname = ?", flowName).
		First(&wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Errorf("workflow %s does not exist", flowName)
			return err
		}
		return errors.New("INVALID_WORKFLOW")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(wf.Flow), &parsed); err != nil {
		fmt.Errorf("failed to decode flow JSON: %w", err)
		return err
	}

	df, ok := parsed["drawflow"].(map[string]interface{})
	if !ok {
		return errors.New("INVALID_WORKFLOW_STRUCTURE")
	}

	home, ok := df["Home"].(map[string]interface{})
	if !ok {
		return errors.New("INVALID_WORKFLOW_STRUCTURE")
	}

	dataFlow, ok := home["data"].(map[string]interface{})
	if !ok {
		return errors.New("INVALID_WORKFLOW_STRUCTURE")
	}

	for _, node := range dataFlow {
		nodeMap, _ := node.(map[string]interface{})
		component := Component{
			ID:      fmt.Sprintf("%v", nodeMap["id"]),
			Name:    fmt.Sprintf("%v", nodeMap["name"]),
			Outputs: make(map[string]interface{}),
		}
		if outputs, ok := nodeMap["outputs"].(map[string]interface{}); ok {
			component.Outputs = outputs
		}
		components = append(components, component)
	}

	data, _ := json.MarshalIndent(components, "", "  ")
	fmt.Println(string(data))

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

	if len(components) > 0 {
		componentId, _ := strconv.Atoi(components[0].ID)
		if err := InternalFlow(c, components[0], int(wf.Workflowid), componentId, db, search); err != nil {
			if !search {
				tx.Rollback()
			}
			fmt.Errorf("internal flow failed: %w", err)
			return err
		}
	}

	if !search {
		if err := tx.Commit().Error; err != nil {
			fmt.Errorf("commit failed: %w", err)
			return err
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
