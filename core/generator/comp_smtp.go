package generator

import (
	"crypto/tls"
	"erp6-be-golang/core/helpers"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("smtp", func(ctx *WorkflowContext) error {
		return handleSMTP(ctx.FiberCtx, ctx.Params, ctx.DB)
	})
}

// Template represents the template table structure
type Template struct {
	TemplateID   int    `gorm:"column:templateid;primaryKey"`
	TemplateName string `gorm:"column:templatename"`
	Content      string `gorm:"column:content"`
	CreatedBy    int    `gorm:"column:createdby"`
	CreatedDate  string `gorm:"column:createddate"`
	UpdatedBy    int    `gorm:"column:updatedby"`
	UpdatedDate  string `gorm:"column:updateddate"`
}

func (Template) TableName() string {
	return "template"
}

func handleSMTP(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB) error {
	var (
		host         string
		port         int
		username     string
		password     string
		from         string
		to           string
		subject      string
		message      string
		templateID   int
		templateName string
	)

	// Extract parameters from workflow
	for _, p := range params {
		val := strings.TrimSpace(ResolveParam(c, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "host":
			host = val
		case "port":
			if val != "" {
				fmt.Sscanf(val, "%d", &port)
			}
		case "username":
			username = val
		case "password":
			password = val
		case "from":
			from = val
		case "to":
			to = val
		case "subject":
			subject = val
		case "message":
			message = val
		case "templateid":
			if val != "" {
				fmt.Sscanf(val, "%d", &templateID)
			}
		case "templatename":
			templateName = val
		}
	}

	// Fallback to environment variables for SMTP configuration
	if host == "" {
		host = os.Getenv("SMTP_HOST")
	}
	if port == 0 {
		if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
			fmt.Sscanf(portStr, "%d", &port)
		}
	}
	if username == "" {
		username = os.Getenv("SMTP_USERNAME")
	}
	if password == "" {
		password = os.Getenv("SMTP_PASSWORD")
	}
	if from == "" {
		from = os.Getenv("SMTP_FROM")
	}

	// Load template from database if templateid or templatename is provided
	if templateID > 0 || templateName != "" {
		var tmpl Template
		var err error

		if templateID > 0 {
			err = db.Where("templateid = ?", templateID).First(&tmpl).Error
		} else {
			err = db.Where("templatename = ?", templateName).First(&tmpl).Error
		}

		if err == nil {
			// Template found, use its content
			message = tmpl.Content
			log.Infof("[SMTP] Loaded template: %s (ID: %d)", tmpl.TemplateName, tmpl.TemplateID)
		} else if err == gorm.ErrRecordNotFound {
			// Template not found, log warning but continue with message parameter
			log.Warnf("[SMTP] Template not found (ID: %d, Name: %s), using message parameter", templateID, templateName)
		} else {
			// Database error
			helpers.FailResponse(c, fiber.StatusInternalServerError, "TEMPLATE_LOAD_ERROR", err.Error())
			return nil
		}
	}

	// Resolve variables in message content
	if message != "" {
		message = ResolveParam(c, message)
	}

	// Validate required parameters
	if host == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "SMTP_ERROR", "host is required")
		return nil
	}
	if port == 0 {
		helpers.FailResponse(c, fiber.StatusBadRequest, "SMTP_ERROR", "port is required")
		return nil
	}
	if from == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "SMTP_ERROR", "from is required")
		return nil
	}
	if to == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "SMTP_ERROR", "to is required")
		return nil
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", from)

	// Handle multiple recipients
	toAddresses := strings.Split(to, ",")
	trimmedTo := make([]string, 0, len(toAddresses))
	for _, addr := range toAddresses {
		trimmed := strings.TrimSpace(addr)
		if trimmed != "" {
			trimmedTo = append(trimmedTo, trimmed)
		}
	}
	m.SetHeader("To", trimmedTo...)

	if subject != "" {
		m.SetHeader("Subject", subject)
	}

	if message != "" {
		m.SetBody("text/html", message)
	}

	// Create SMTP dialer
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Errorf("[SMTP] Failed to send email: %v", err)
		helpers.FailResponse(c, fiber.StatusInternalServerError, "SMTP_SEND_ERROR", err.Error())
		return nil
	}

	log.Infof("[SMTP] Email sent successfully to: %s", to)

	// Prepare result
	resultStat := map[string]interface{}{
		"status":  "success",
		"message": "Email sent successfully",
		"to":      to,
		"subject": subject,
	}

	// Inject result into workflow context
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{
		DataInputNode: "",
		ResultNode:    resultStat,
	})
	c.Locals("wfEngine", wfEngine)

	helpers.SuccessResponse(c, "EMAIL_SENT", resultStat)
	return nil
}
