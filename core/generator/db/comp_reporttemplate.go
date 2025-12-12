package generator

import (
	"erp6-be-golang/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("reporttemplate", func(ctx *WorkflowContext) error {
		return handleReportTemplate(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.Search)
	})
}

func handleReportTemplate(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, search bool) error {
	var (
		templateID   string
		outputFormat string
		parameters   = make(map[string]string)
	)

	// Parse workflow parameters
	for _, v := range params {
		switch v.InputName {
		case "templateid":
			templateID = strings.TrimSpace(v.CompValue)
		case "format":
			outputFormat = strings.TrimSpace(v.CompValue)
		default:
			// All other parameters are report parameters
			if v.InputName != "" {
				parameters[v.InputName] = GetSearchText(c, []string{"POST"}, v.InputName, "", "string")
			}
		}
	}

	if templateID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Template ID is required",
		})
	}

	// Get report template
	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ? AND recordstatus = 1", templateID).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    404,
				"message": "Report template not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to fetch report template",
			"error":   err.Error(),
		})
	}

	// Use template's default format if not specified
	if outputFormat == "" {
		outputFormat = strings.ToLower(template.ReportType)
	}

	// Convert parameters to interface map
	paramMap := make(map[string]interface{})
	for k, v := range parameters {
		paramMap[k] = v
	}

	// Store result in workflow context for next components
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	resultData := map[string]interface{}{
		"templateid": template.ReportTemplateID,
		"reportname": template.ReportName,
		"format":     outputFormat,
		"parameters": paramMap,
	}

	wfEngine = append(wfEngine, WorkflowEngine{
		DataInputNode: "",
		ResultNode:    resultData,
	})
	c.Locals("wfEngine", wfEngine)

	// Return success response
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Report template processed successfully",
		"data":    resultData,
	})
}
