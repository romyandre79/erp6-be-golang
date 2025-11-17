package admin

import (
	"erp6-be-golang/core/events"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitMenuaccess() {
	events.Register("BeforeGet:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:menuaccess", func(data interface{}) error {
		return nil
	})
}

// menuSingleNameHandler godoc
func MenuSingleNameHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid")
	menuName := c.Query("menuname")

	if userID == nil {
		return helpers.FailResponse(c, fiber.StatusNotFound, "INVALID_USER", "NO_USER_FOUND")
	}

	// --- STEP 4: Ambil hanya data menu yang dibolehkan ---
	var menus models.Menuaccess
	err := db.
		Table("menuaccess").
		Preload("Modules").
		Where("menuaccess.menuname = ?", menuName).
		Find(&menus).Error
	if err != nil {
		return helpers.FailResponse(c, fiber.StatusInternalServerError, "MENU_QUERY_FAILED", err.Error())
	}

	return helpers.SuccessResponse(c, "DATA RETRIEVED", menus)
}
