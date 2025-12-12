package admin

import (
	"encoding/json"
	"erp6-be-golang/core/jrxml"
	"erp6-be-golang/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateReportTemplate creates a new report template
func CreateReportTemplate(c *fiber.Ctx, db *gorm.DB) error {
	var template models.ReportTemplate

	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	template.UpdateDate = time.Now()
	template.RecordStatus = 1

	if err := db.Create(&template).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to create report template",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Report template created successfully",
		"data":    template,
	})
}

// ListReportTemplates lists all report templates
func ListReportTemplates(c *fiber.Ctx, db *gorm.DB) error {
	var templates []models.ReportTemplate

	query := db.Where("recordstatus = 1")

	// Optional filters
	if category := c.Query("category"); category != "" {
		query = query.Where("reportcategory = ?", category)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("reportname LIKE ? OR reportdesc LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&templates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to fetch report templates",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success",
		"data":    templates,
	})
}

// GetReportTemplate gets a single report template by ID
func GetReportTemplate(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ? AND recordstatus = 1", id).First(&template).Error; err != nil {
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

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success",
		"data":    template,
	})
}

// UpdateReportTemplate updates an existing report template
func UpdateReportTemplate(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ?", id).First(&template).Error; err != nil {
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

	var updates models.ReportTemplate
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	updates.UpdateDate = time.Now()
	updates.ReportTemplateID = template.ReportTemplateID

	if err := db.Model(&template).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to update report template",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Report template updated successfully",
		"data":    template,
	})
}

// DeleteReportTemplate soft deletes a report template
func DeleteReportTemplate(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	result := db.Model(&models.ReportTemplate{}).
		Where("reporttemplateid = ?", id).
		Update("recordstatus", 0)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to delete report template",
			"error":   result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "Report template not found",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Report template deleted successfully",
	})
}

// ImportJRXML imports a JRXML file and converts it to our format
func ImportJRXML(c *fiber.Ctx, db *gorm.DB) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to read file",
			"error":   err.Error(),
		})
	}
	defer fileContent.Close()

	// Read file bytes
	fileBytes := make([]byte, file.Size)
	if _, err := fileContent.Read(fileBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to read file content",
			"error":   err.Error(),
		})
	}

	jrxmlContent := string(fileBytes)

	// Convert JRXML to JSON
	jsonContent, err := jrxml.ConvertJRXMLToJSON(jrxmlContent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Failed to parse JRXML file",
			"error":   err.Error(),
		})
	}

	// Parse JSON to get page dimensions
	var templateData jrxml.ReportTemplate
	if err := json.Unmarshal([]byte(jsonContent), &templateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to parse converted JSON",
			"error":   err.Error(),
		})
	}

	// Create new report template
	template := models.ReportTemplate{
		ReportName:     c.FormValue("reportname"),
		ReportDesc:     c.FormValue("reportdesc"),
		ReportCategory: c.FormValue("reportcategory"),
		ReportType:     "PDF",
		TemplateJSON:   jsonContent,
		JRXMLContent:   jrxmlContent,
		PageWidth:      templateData.PageWidth,
		PageHeight:     templateData.PageHeight,
		Orientation:    templateData.Orientation,
		RecordStatus:   1,
		UpdateDate:     time.Now(),
	}

	if moduleID := c.FormValue("moduleid"); moduleID != "" {
		if id, err := strconv.Atoi(moduleID); err == nil {
			template.ModuleID = id
		}
	}

	if err := db.Create(&template).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to save report template",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "JRXML imported successfully",
		"data":    template,
	})
}

// ExportJRXML exports a report template as JRXML
func ExportJRXML(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ? AND recordstatus = 1", id).First(&template).Error; err != nil {
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

	var jrxmlContent string
	var err error

	// If JRXML content already exists, use it
	if template.JRXMLContent != "" {
		jrxmlContent = template.JRXMLContent
	} else if template.TemplateJSON != "" {
		// Convert JSON to JRXML
		jrxmlContent, err = jrxml.ConvertJSONToJRXML(template.TemplateJSON)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    500,
				"message": "Failed to convert template to JRXML",
				"error":   err.Error(),
			})
		}

		// Save the generated JRXML for future use
		db.Model(&template).Update("jrxmlcontent", jrxmlContent)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "Template has no content to export",
		})
	}

	// Set headers for file download
	filename := fmt.Sprintf("%s.jrxml", template.ReportName)
	c.Set("Content-Type", "application/xml")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Content-Length", strconv.Itoa(len(jrxmlContent)))

	return c.SendString(jrxmlContent)
}

// PreviewReport generates a preview of the report
func PreviewReport(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ? AND recordstatus = 1", id).First(&template).Error; err != nil {
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

	// Parse parameters from request body
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		requestBody = make(map[string]interface{})
	}

	// Execute report with preview format (PDF)
	return executeReportInternal(c, db, &template, requestBody, "pdf", true)
}

// ExecuteReport executes a report and returns the generated file
func ExecuteReport(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var template models.ReportTemplate
	if err := db.Where("reporttemplateid = ? AND recordstatus = 1", id).First(&template).Error; err != nil {
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

	// Parse parameters and format from request
	var requestBody struct {
		Parameters map[string]interface{} `json:"parameters"`
		Format     string                 `json:"format"` // pdf, xls, csv
	}
	if err := c.BodyParser(&requestBody); err != nil {
		requestBody.Parameters = make(map[string]interface{})
		requestBody.Format = "pdf"
	}

	if requestBody.Format == "" {
		requestBody.Format = strings.ToLower(template.ReportType)
	}

	return executeReportInternal(c, db, &template, requestBody.Parameters, requestBody.Format, false)
}

// executeReportInternal handles the actual report execution
func executeReportInternal(c *fiber.Ctx, db *gorm.DB, template *models.ReportTemplate, parameters map[string]interface{}, format string, isPreview bool) error {
	// Get JasperReports server configuration from environment
	reportURL := os.Getenv("REPORT_URL")
	reportUser := os.Getenv("REPORT_USER")
	reportPass := os.Getenv("REPORT_PASS")

	if reportURL == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"code":    503,
			"message": "JasperReports server not configured",
		})
	}

	// If template has JRXML content, we need to upload it to JasperReports server first
	// For now, we'll use the existing report server infrastructure
	// This assumes the JRXML is already deployed to the JasperReports server

	// Build query parameters
	queryParams := url.Values{}
	for key, value := range parameters {
		queryParams.Add(key, fmt.Sprintf("%v", value))
	}

	// Determine output format
	outputFormat := format
	switch strings.ToLower(format) {
	case "pdf":
		outputFormat = "pdf"
	case "xls", "xlsx", "excel":
		outputFormat = "xlsx"
	case "csv":
		outputFormat = "csv"
	default:
		outputFormat = "pdf"
	}

	// Construct JasperReports server URL
	// Note: This is a simplified implementation
	// You may need to adjust based on your JasperReports server setup
	fullURL := fmt.Sprintf("%s.%s?%s", reportURL, outputFormat, queryParams.Encode())

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to create request",
			"error":   err.Error(),
		})
	}

	// Add basic authentication
	if reportUser != "" && reportPass != "" {
		req.SetBasicAuth(reportUser, reportPass)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"code":    503,
			"message": "Failed to connect to JasperReports server",
			"error":   err.Error(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"code":    resp.StatusCode,
			"message": "JasperReports server returned error",
		})
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "Failed to read response",
			"error":   err.Error(),
		})
	}

	// Set appropriate content type and headers
	var contentType string
	var fileExtension string
	switch outputFormat {
	case "pdf":
		contentType = "application/pdf"
		fileExtension = "pdf"
	case "xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		fileExtension = "xlsx"
	case "csv":
		contentType = "text/csv"
		fileExtension = "csv"
	default:
		contentType = "application/octet-stream"
		fileExtension = "bin"
	}

	c.Set("Content-Type", contentType)

	if !isPreview {
		filename := fmt.Sprintf("%s.%s", template.ReportName, fileExtension)
		c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	}

	c.Set("Content-Length", strconv.Itoa(len(body)))

	return c.Send(body)
}
