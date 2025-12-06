package models

import "time"

type Widget struct {
	Widgetid      int       `gorm:"column:widgetid;primaryKey;autoIncrement" json:"widgetid"`
	Widgetname    string    `gorm:"column:widgetname" json:"widgetname"`
	Widgettitle   string    `gorm:"column:widgettitle" json:"widgettitle"`
	Widgetversion string    `gorm:"column:widgetversion" json:"widgetversion"`
	Widgetby      string    `gorm:"column:widgetby" json:"widgetby"`
	Description   string    `gorm:"column:description" json:"description"`
	Widgetform    string    `gorm:"column:widgetform" json:"widgetform"`
	Moduleid      int       `gorm:"column:moduleid" json:"moduleid"`
	Recordstatus  int       `gorm:"column:recordstatus;default:1" json:"recordstatus"`
	Createdate    time.Time `gorm:"column:createdate;autoCreateTime" json:"createdate"`
	Updatedate    time.Time `gorm:"column:updatedate;autoUpdateTime" json:"updatedate"`
}

func (Widget) TableName() string {
	return "widget"
}
