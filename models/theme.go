package models

import "time"

type Theme struct {
	Themeid      int       `gorm:"column:themeid;primaryKey;autoIncrement"`
	Themename    string    `gorm:"column:themename;size:20;unique;not null"`
	Description  string    `gorm:"column:description;type:text;not null"`
	Createdby    string    `gorm:"column:createdby;size:30;not null"`
	Themeversion string    `gorm:"column:themeversion;size:3;not null"`
	Themedata    string    `gorm:"column:themedata;type:text;not null"`
	Installdate  time.Time `gorm:"column:installdate;autoCreateTime"`
	Recordstatus int8      `gorm:"column:recordstatus;default:0;not null"`
	Updatedate   time.Time `gorm:"column:updatedate;autoUpdateTime"`
}

func (Theme) TableName() string {
	return "theme"
}
