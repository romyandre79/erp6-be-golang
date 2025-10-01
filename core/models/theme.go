package models

import "time"

type Theme struct {
	Themeid      int    `gorm:"column:themeid;primaryKey;autoIncrement"`
	Themename    string `gorm:"column:themename"`
	Description  string
	Isadmin      int
	Createdby    string
	Themeversion string
	Installdate  time.Time
	Recordstatus int
	Updatedate   time.Time
}

func (Theme) TableName() string {
	return "theme"
}
