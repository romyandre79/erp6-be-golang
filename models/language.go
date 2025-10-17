package models

import "time"

type Language struct {
	Languageid int `gorm:"column:languageid;primaryKey" json:"languageid"`
	Languagename string `gorm:"column:languagename" json:"languagename"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Language) TableName() string {
	return "language"
}
