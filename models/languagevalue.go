package models

import (
  "time"
)

type Languagevalue struct {
	Languagevalid int `gorm:"column:languagevalid;primaryKey" json:"languagevalid"`
	Languagevalname string `gorm:"column:languagevalname" json:"languagevalname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Languagevalue) TableName() string {
	return "languagevalue"
}
