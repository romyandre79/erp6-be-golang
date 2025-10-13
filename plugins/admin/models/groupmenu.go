package admin

import (
	"erp6-be-golang/core/models"
	"time"
)

type GroupMenu struct {
	GroupMenuID   int `gorm:"column:groupmenuid;primaryKey;autoIncrement" json:"groupmenuid"`
	GroupAccessID int `gorm:"column:groupaccessid;not null" json:"groupaccessid"`
	MenuAccessID  int `gorm:"column:menuaccessid;not null" json:"menuaccessid"`

	IsRead     int `gorm:"column:isread;default:1" json:"isread"`
	IsWrite    int `gorm:"column:iswrite;default:1" json:"iswrite"`
	IsPost     int `gorm:"column:ispost;default:1" json:"ispost"`
	IsReject   int `gorm:"column:isreject;default:1" json:"isreject"`
	IsUpload   int `gorm:"column:isupload;default:1" json:"isupload"`
	IsDownload int `gorm:"column:isdownload;default:1" json:"isdownload"`
	IsPurge    int `gorm:"column:ispurge;default:1" json:"ispurge"`

	UpdatedAt time.Time `gorm:"column:updatedate;autoUpdateTime" json:"updatedate"`

	// Relasi
	GroupAccess GroupAccess       `gorm:"foreignKey:GroupAccessID;references:GroupAccessID" json:"group_access"`
	MenuAccess  models.MenuAccess `gorm:"foreignKey:MenuAccessID;references:MenuAccessID" json:"menu_access"`
}

// TableName override
func (GroupMenu) TableName() string {
	return "groupmenu"
}
