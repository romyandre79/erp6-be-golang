package models

import (
	"time"
)

type Widget struct {
	Widgetid      int       `gorm:"column:widgetid;primaryKey" json:"widgetid"`
	Widgetname    string    `gorm:"column:widgetname" json:"widgetname"`
	Widgettitle   string    `gorm:"column:widgettitle" json:"widgettitle"`
	Widgetversion string    `gorm:"column:widgetversion" json:"widgetversion"`
	Widgetby      string    `gorm:"column:widgetby" json:"widgetby"`
	Description   string    `gorm:"column:description" json:"description"`
	Widgeturl     string    `gorm:"column:widgeturl" json:"widgeturl"`
	Moduleid      int       `gorm:"column:moduleid" json:"moduleid"`
	Installdate   time.Time `gorm:"column:installdate" json:"installdate"`
	Recordstatus  int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`

	Modules *Modules `gorm:"foreignKey:Moduleid;references:Moduleid" json:"modules"`
}

func (Widget) TableName() string {
	return "widget"
}
