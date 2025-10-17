package admin

import (
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PermissionAction string

const (
	PermRead     PermissionAction = "read"
	PermWrite    PermissionAction = "write"
	PermPost     PermissionAction = "post"
	PermReject   PermissionAction = "reject"
	PermUpload   PermissionAction = "upload"
	PermDownload PermissionAction = "download"
	PermPurge    PermissionAction = "purge"
)

// check User Permission based on menu name and action
func CheckUserPermission(c *fiber.Ctx, db *gorm.DB, menuName string, action PermissionAction) (bool, error) {
	userID := c.Locals("userid")

	if userID == nil {
		return false, helpers.FailResponse(c, fiber.StatusUnauthorized, "INVALID_USER", "")
	}

	var user models.Useraccess
	if err := db.Preload("Usergroups.Groupaccess.Groupmenus.Menuaccess").
		First(&user, userID).Error; err != nil {
		return false, err
	}

	for _, ug := range user.Usergroups {
		for _, gm := range ug.Groupaccess.Groupmenus {
			if gm.Menuaccess.Menuname == menuName {
				switch action {
				case PermRead:
					return gm.Isread == 1, nil
				case PermWrite:
					return gm.Iswrite == 1, nil
				case PermPost:
					return gm.Ispost == 1, nil
				case PermReject:
					return gm.Isreject == 1, nil
				case PermUpload:
					return gm.Isupload == 1, nil
				case PermDownload:
					return gm.Isdownload == 1, nil
				case PermPurge:
					return gm.Ispurge == 1, nil
				default:
					return false, errors.New("invalid action")
				}
			}
		}
	}
	return false, nil
}
