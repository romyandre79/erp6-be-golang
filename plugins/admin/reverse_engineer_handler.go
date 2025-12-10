package admin

import (
	"fmt"
	"time"

	generator "erp6-be-golang/core/generator/db"
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ReverseEngineerRequest represents the request body for reverse engineering
type ReverseEngineerRequest struct {
	AutoLayout bool `json:"auto_layout"` // Whether to auto-layout tables in a grid
}

// ReverseEngineerResponse represents the response from reverse engineering
type ReverseEngineerResponse struct {
	Success       bool     `json:"success"`
	Message       string   `json:"message"`
	TablesCreated int      `json:"tables_created"`
	TablesUpdated int      `json:"tables_updated"`
	TableNames    []string `json:"table_names"`
	Error         string   `json:"error,omitempty"`
}

// ReverseEngineerHandler handles the reverse engineering request
func ReverseEngineerHandler(c *fiber.Ctx, db *gorm.DB) error {
	var req ReverseEngineerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ReverseEngineerResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid request body: %v", err),
		})
	}

	// Default to auto-layout if not specified
	if !req.AutoLayout {
		req.AutoLayout = true
	}

	// Extract all tables from database
	tables, err := generator.ReverseEngineerDatabase(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ReverseEngineerResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to reverse engineer database: %v", err),
		})
	}

	if len(tables) == 0 {
		return c.JSON(ReverseEngineerResponse{
			Success: true,
			Message: "No tables found in database",
		})
	}

	// Layout configuration
	const spacingX = 340
	const spacingY = 200
	const perRow = 3
	startX := 40
	startY := 40

	var createdCount, updatedCount int
	var tableNames []string

	// Process each table
	for i, table := range tables {
		// Calculate position for auto-layout
		x := startX
		y := startY
		if req.AutoLayout {
			row := i / perRow
			col := i % perRow
			x = startX + col*spacingX
			y = startY + row*spacingY
		}

		// Convert to TableDefinition
		tableDef := generator.ConvertToTableDefinition(&table, x, y)

		// Convert to JSON
		jsonContent, err := generator.ConvertToJSON(tableDef)
		if err != nil {
			fmt.Printf("Warning: failed to convert table %s to JSON: %v\n", table.Name, err)
			continue
		}

		// Check if dbobject already exists
		var existingDbObject models.Dbobject
		result := db.Where("objectname = ? AND objecttype = ?", table.Name, "table").First(&existingDbObject)

		if result.Error == nil {
			// Update existing
			existingDbObject.Objectcontent = jsonContent
			existingDbObject.Objectversion++
			existingDbObject.Comment = "Reverse engineered from database (updated)"
			existingDbObject.Updatedate = time.Now()

			if err := db.Save(&existingDbObject).Error; err != nil {
				fmt.Printf("Warning: failed to update dbobject for table %s: %v\n", table.Name, err)
				continue
			}
			updatedCount++
		} else if result.Error == gorm.ErrRecordNotFound {
			// Create new
			newDbObject := models.Dbobject{
				Objectname:    table.Name,
				Objecttype:    "table",
				Objectcontent: jsonContent,
				Objectversion: 1,
				Ispublished:   0,
				Sort:          int8(i),
				Comment:       "Reverse engineered from database",
				Updatedate:    time.Now(),
			}

			if err := db.Create(&newDbObject).Error; err != nil {
				fmt.Printf("Warning: failed to create dbobject for table %s: %v\n", table.Name, err)
				continue
			}
			createdCount++
		} else {
			fmt.Printf("Warning: database error checking table %s: %v\n", table.Name, result.Error)
			continue
		}

		tableNames = append(tableNames, table.Name)
	}

	message := fmt.Sprintf("Reverse engineering completed. Created: %d, Updated: %d", createdCount, updatedCount)

	return c.JSON(ReverseEngineerResponse{
		Success:       true,
		Message:       message,
		TablesCreated: createdCount,
		TablesUpdated: updatedCount,
		TableNames:    tableNames,
	})
}
