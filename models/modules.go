package models

import "time"

type Modules struct {
	Moduleid      int       `gorm:"column:moduleid;primaryKey" json:"moduleid"`
	Modulename    string    `gorm:"column:modulename" json:"modulename"`
	Description   string    `gorm:"column:description" json:"description"`
	Createdby     string    `gorm:"column:createdby" json:"createdby"`
	Moduleversion string    `gorm:"column:moduleversion" json:"moduleversion"`
	Installdate   time.Time `gorm:"column:installdate" json:"installdate"`
	Themeid       int       `gorm:"column:themeid" json:"themeid"`
	Recordstatus  int       `gorm:"column:recordstatus" json:"recordstatus"`
	Theme         Theme     `gorm:"foreignKey:themeid;references:themeid"`
}

func (Modules) TableName() string {
	return "modules"
}
