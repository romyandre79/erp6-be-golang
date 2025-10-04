package models

import (
	"time"
)

type MenuAccess struct {
	MenuAccessID   int       `gorm:"column:menuaccessid;primaryKey;autoIncrement"`
	MenuName       string    `gorm:"column:menuname;size:20;not null;uniqueIndex:uq_menuaccess_menuname"`
	MenuCode       string    `gorm:"column:menucode;size:50"`
	Description    string    `gorm:"column:description;size:50;not null"`
	ModuleID       int       `gorm:"column:moduleid;not null"`
	ParentID       *int      `gorm:"column:parentid"`
	MenuURL        *string   `gorm:"column:menuurl;size:30"`
	SortOrder      int       `gorm:"column:sortorder;not null;default:1"`
	MenuIcon       *string   `gorm:"column:menuicon;size:50"`
	MenuType       string    `gorm:"column:menutype;size:50;not null;default:LIST"`
	MenuFormDetail *string   `gorm:"column:menuformdetail;type:mediumtext"`
	RecordStatus   int8      `gorm:"column:recordstatus;not null;default:0"`
	MenuForm       *string   `gorm:"column:menuform;type:mediumtext"`
	UpdateDate     time.Time `gorm:"column:updatedate;autoUpdateTime"`

	// Relasi ke Module
	Module Module `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relasi ke parent menu (self join)
	Parent   *MenuAccess  `gorm:"foreignKey:ParentID"`
	Children []MenuAccess `gorm:"foreignKey:ParentID"`
}

// Nama tabel
func (MenuAccess) TableName() string {
	return "menuaccess"
}
