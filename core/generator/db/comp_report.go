package generator

import (
	"erp6-be-golang/core/configs"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/core/i18n"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("reportserver", func(ctx *WorkflowContext) error {
		return handleReportServer(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
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
