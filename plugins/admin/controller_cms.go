package admin

import (
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PublicPageHandler retrieves a public page schema by slug
func PublicPageHandler(c *fiber.Ctx, db *gorm.DB) error {
	slug := c.Params("*")

	if slug == "" {
		return helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_SLUG", "Slug is required")
	}

	var menu models.Menuaccess
	// Find menu by name or code, MUST be type 'page' for security
	err := db.Where("(menuname = ? OR menucode = ?) AND menutype = ? AND recordstatus = 1", slug, slug, "page").First(&menu).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helpers.FailResponse(c, fiber.StatusNotFound, "PAGE_NOT_FOUND", "Page not found")
		}
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error())
	}

	return helpers.SuccessResponse(c, "PAGE_RETRIEVED", menu)
}
