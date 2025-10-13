package admin

import (
	admin "erp6-be-golang/plugins/admin/models"
	"errors"

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

// Can check apakah user punya izin terhadap menu dengan kode tertentu
func GetGroupRole(db *gorm.DB, userID int, menuCode string, action PermissionAction) (bool, error) {
	var user admin.UserAccess
	if err := db.Preload("Groups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, userID).Error; err != nil {
		return false, err
	}

	for _, ug := range user.Groups {
		for _, gm := range ug.GroupAccess.GroupMenus {
			if gm.MenuAccess.MenuCode == menuCode {
				switch action {
				case PermRead:
					return gm.IsRead == 1, nil
				case PermWrite:
					return gm.IsWrite == 1, nil
				case PermPost:
					return gm.IsPost == 1, nil
				case PermReject:
					return gm.IsReject == 1, nil
				case PermUpload:
					return gm.IsUpload == 1, nil
				case PermDownload:
					return gm.IsDownload == 1, nil
				case PermPurge:
					return gm.IsPurge == 1, nil
				default:
					return false, errors.New("invalid action")
				}
			}
		}
	}
	return false, nil
}

func LoadUserWithPermissions(db *gorm.DB, userID int) (*admin.UserAccess, error) {
	var user admin.UserAccess

	err := db.Preload("Groups.GroupAccess.GroupMenus.MenuAccess").
		First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
