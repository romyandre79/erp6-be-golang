package admin

import (
	"erp6-be-golang/core/models"
	"time"
)

type GroupMenu struct {
	GroupMenuID   int `gorm:"column:groupmenuid;primaryKey;autoIncrement" json:"groupmenuid"`
	GroupAccessID int `gorm:"column:groupaccessid;not null" json:"groupaccessid"`
	MenuAccessID  int `gorm:"column:menuaccessid;not null" json:"menuaccessid"`

	IsRead     bool `gorm:"column:isread;default:1" json:"isread"`
	IsWrite    bool `gorm:"column:iswrite;default:1" json:"iswrite"`
	IsPost     bool `gorm:"column:ispost;default:1" json:"ispost"`
	IsReject   bool `gorm:"column:isreject;default:1" json:"isreject"`
	IsUpload   bool `gorm:"column:isupload;default:1" json:"isupload"`
	IsDownload bool `gorm:"column:isdownload;default:1" json:"isdownload"`
	IsPurge    bool `gorm:"column:ispurge;default:1" json:"ispurge"`

	UpdatedAt time.Time `gorm:"column:updatedate;autoUpdateTime" json:"updatedate"`

	// Relasi
	GroupAccess GroupAccess       `gorm:"foreignKey:GroupAccessID;references:GroupAccessID" json:"group_access"`
	MenuAccess  models.MenuAccess `gorm:"foreignKey:MenuAccessID;references:MenuAccessID" json:"menu_access"`
}

// TableName override
func (GroupMenu) TableName() string {
	return "groupmenu"
}
