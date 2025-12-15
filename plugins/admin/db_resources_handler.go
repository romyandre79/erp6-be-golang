package admin

import (
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Helper to save relations
func SaveRelationsHandler(c *fiber.Ctx, db *gorm.DB) error {
	var input []models.DbobjectRelation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Simple replacement strategy: delete all and re-create (or upsert if you prefer)
	// For a designer, typically we might want to replace the entire set or handle individually.
	// Assuming bulk save for the entire canvas/project usually implies a "state save".
	// But since this is a global table, we probably want to filter by project/diagram ID if that existed.
	// Given the current simple schema, we'll assume we are saving ALL active relations or passing a list to upset.

	// WARNING: This implementation assumes we are saving specific relations.
	// If the frontend sends the "entire list of relations to keep", we should potentially clear old ones.
	// However, without a project ID, clearing "all" is dangerous if used by multiple things.

	// For now, let's implement UPSERT based on ID.
	for i := range input {
		if input[i].ID > 0 {
			db.Save(&input[i])
		} else {
			// Check for existing relation to avoid duplicates
			var existing models.DbobjectRelation
			err := db.Where("fromtableid = ? AND fromcolindex = ? AND totableid = ? AND tocolindex = ?",
				input[i].FromTableID, input[i].FromColIndex, input[i].ToTableID, input[i].ToColIndex).First(&existing).Error

			if err == nil {
				// Found existing, update it
				input[i].ID = existing.ID
				db.Save(&input[i])
			} else {
				// Not found, create new
				db.Create(&input[i])
			}
		}
	}

	return c.JSON(fiber.Map{"success": true, "data": input})
}

// Helper to get relations
func GetRelationsHandler(c *fiber.Ctx, db *gorm.DB) error {
	var relations []models.DbobjectRelation
	db.Find(&relations)
	return c.JSON(relations)
}

// Helper to delete relations
func DeleteRelationHandler(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	db.Delete(&models.DbobjectRelation{}, id)
	return c.JSON(fiber.Map{"success": true})
}

// Same for Areas
func SaveAreasHandler(c *fiber.Ctx, db *gorm.DB) error {
	var input []models.DbobjectArea
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	for i := range input {
		if input[i].ID > 0 {
			db.Save(&input[i])
		} else {
			db.Create(&input[i])
		}
	}

	return c.JSON(fiber.Map{"success": true, "data": input})
}

func GetAreasHandler(c *fiber.Ctx, db *gorm.DB) error {
	var areas []models.DbobjectArea
	db.Find(&areas)
	return c.JSON(areas)
}

func DeleteAreaHandler(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	db.Delete(&models.DbobjectArea{}, id)
	return c.JSON(fiber.Map{"success": true})
}
