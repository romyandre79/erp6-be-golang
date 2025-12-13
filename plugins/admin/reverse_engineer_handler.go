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
		fmt.Printf("DEBUG: Processing extracted table: %s\n", table.Name)
		
		// Default position
		x := startX
		y := startY
		if req.AutoLayout {
			row := i / perRow
			col := i % perRow
			x = startX + col*spacingX
			y = startY + row*spacingY
		}

		// Check if dbobject already exists to preserve coordinates
		var existingDbObject models.Dbobject
		result := db.Where("objectname = ? AND objecttype = ?", table.Name, "table").First(&existingDbObject)
		isExisting := result.Error == nil

		if isExisting {
			fmt.Printf("DEBUG: Found existing dbobject for table: %s. Preserving coordinates.\n", table.Name)
			// Try to parse existing content to get X and Y
			existingDef, err := generator.ParseTableJSON(existingDbObject.Objectcontent)
			if err == nil && existingDef != nil {
				x = existingDef.Table.X
				y = existingDef.Table.Y
				fmt.Printf("DEBUG: Preserved coordinates for %s: X=%d, Y=%d\n", table.Name, x, y)
			} else {
				fmt.Printf("DEBUG: Failed to parse existing content for %s, using auto-layout coordinates: %v\n", table.Name, err)
			}
		} else {
			fmt.Printf("DEBUG: New table detected: %s\n", table.Name)
		}

		// Convert to TableDefinition with (potentially preserved) coordinates
		tableDef := generator.ConvertToTableDefinition(&table, x, y)

		// Convert to JSON
		jsonContent, err := generator.ConvertToJSON(tableDef)
		if err != nil {
			fmt.Printf("Warning: failed to convert table %s to JSON: %v\n", table.Name, err)
			continue
		}

		if isExisting {
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
			fmt.Printf("DEBUG: Updated table %s\n", table.Name)
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
			fmt.Printf("DEBUG: Created table %s\n", table.Name)
		} else {
			fmt.Printf("Warning: database error checking table %s: %v\n", table.Name, result.Error)
			continue
		}

		tableNames = append(tableNames, table.Name)
	}

	// Process Relationships (Foreign Keys)
	var relationsCreated int
	fmt.Println("DEBUG: Starting relationship extraction...")

	// Need to reload table map with IDs since we might have created new ones
	tableIDMap := make(map[string]int)
	var allDbObjects []models.Dbobject
	db.Where("objecttype = ?", "table").Find(&allDbObjects)
	for _, obj := range allDbObjects {
		tableIDMap[obj.Objectname] = obj.Dbobjectid
	}

	for _, table := range tables {
		for _, fk := range table.ForeignKeys {
			fromTableID, okFrom := tableIDMap[table.Name]
			toTableID, okTo := tableIDMap[fk.ReferencedTable]

			if !okFrom || !okTo {
				fmt.Printf("Warning: Skipping relation for %s -> %s (One or both tables not found in map)\n", table.Name, fk.ReferencedTable)
				continue
			}

			// Find column indices
			fromColIdx := -1
			for idx, col := range table.Columns {
				if col.Name == fk.ColumnName {
					fromColIdx = idx
					break
				}
			}

			// For referenced column index, we need to inspect the referenced table struct (which we have in `tables`)
			toColIdx := -1
			for _, refTable := range tables {
				if refTable.Name == fk.ReferencedTable {
					for idx, col := range refTable.Columns {
						if col.Name == fk.ReferencedColumn {
							toColIdx = idx
							break
						}
					}
					break
				}
			}

			if fromColIdx == -1 || toColIdx == -1 {
				fmt.Printf("Warning: Skipping relation for %s.%s -> %s.%s (Column index not found)\n", 
					table.Name, fk.ColumnName, fk.ReferencedTable, fk.ReferencedColumn)
				continue
			}

			// Check if relation already exists
			var existingRel models.DbobjectRelation
			err := db.Where("fromtableid = ? AND fromcolname = ? AND totableid = ? AND tocolname = ?", 
				fromTableID, fk.ColumnName, toTableID, fk.ReferencedColumn).First(&existingRel).Error

			if err == gorm.ErrRecordNotFound {
				// Create new relation
				newRel := models.DbobjectRelation{
					FromTableID:  fromTableID,
					FromColIndex: fromColIdx,
					FromColName:  fk.ColumnName,
					ToTableID:    toTableID,
					ToColIndex:   toColIdx,
					ToColName:    fk.ReferencedColumn,
					Path:         "", // Auto-routing will handle this in frontend
				}
				if err := db.Create(&newRel).Error; err != nil {
					fmt.Printf("Warning: Failed to create relation: %v\n", err)
				} else {
					relationsCreated++
					fmt.Printf("DEBUG: Created relation %s.%s -> %s.%s\n", table.Name, fk.ColumnName, fk.ReferencedTable, fk.ReferencedColumn)
				}
			} else {
				fmt.Printf("DEBUG: Relation already exists: %s.%s -> %s.%s\n", table.Name, fk.ColumnName, fk.ReferencedTable, fk.ReferencedColumn)
			}
		}
	}

	message := fmt.Sprintf("Reverse engineering completed. Tables: Created %d, Updated %d. Relations: Created %d", createdCount, updatedCount, relationsCreated)

	return c.JSON(ReverseEngineerResponse{
		Success:       true,
		Message:       message,
		TablesCreated: createdCount,
		TablesUpdated: updatedCount,
		TableNames:    tableNames,
	})
}
