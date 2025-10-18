package response

import "time"

type MenuResponse struct {
	MenuAccessID   int       `gorm:"column:menuaccessid;primaryKey;autoIncrement"`
	MenuName       string    `gorm:"column:menuname;size:20;not null;uniqueIndex:uq_menuaccess_menuname"`
	MenuCode       string    `gorm:"column:menucode;size:50"`
	Description    string    `gorm:"column:description;size:50;not null"`
	ModuleID       int       `gorm:"column:moduleid;not null"`
	ParentID       int       `gorm:"column:parentid"`
	ModuleName     string    `gorm:"column:moduleid;not null"`
	ParentName     string    `gorm:"column:parentid"`
	MenuURL        string    `gorm:"column:menuurl;size:30"`
	SortOrder      int       `gorm:"column:sortorder;not null;default:1"`
	MenuIcon       string    `gorm:"column:menuicon;size:50"`
	MenuType       string    `gorm:"column:menutype;size:50;not null;default:LIST"`
	MenuFormDetail string    `gorm:"column:menuformdetail;type:mediumtext"`
	RecordStatus   int8      `gorm:"column:recordstatus;not null;default:0"`
	MenuForm       string    `gorm:"column:menuform;type:mediumtext"`
	UpdateDate     time.Time `gorm:"column:updatedate;autoUpdateTime"`
}
