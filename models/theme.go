package models

import "time"

type Theme struct {
	Themeid int `gorm:"column:themeid;primaryKey" json:"themeid"`
	Themename string `gorm:"column:themename" json:"themename"`
	Description string `gorm:"column:description" json:"description"`
	Isadmin int `gorm:"column:isadmin" json:"isadmin"`
	Createdby string `gorm:"column:createdby" json:"createdby"`
	Themeversion string `gorm:"column:themeversion" json:"themeversion"`
	Installdate time.Time `gorm:"column:installdate" json:"installdate"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Theme) TableName() string {
	return "theme"
}
