package models

import (
  "time"
)

type Employeeforeignlanguage struct {
	Employeeforeignlanguageid int `gorm:"column:employeeforeignlanguageid;primaryKey" json:"employeeforeignlanguageid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Languageid int `gorm:"column:languageid" json:"languageid"`
	Languagevalueid int `gorm:"column:languagevalueid" json:"languagevalueid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Language Language `gorm:"foreignKey:languageid;references:languageid"`
	Languagevalue Languagevalue `gorm:"foreignKey:languagevalueid;references:languagevalid"`
}

func (Employeeforeignlanguage) TableName() string {
	return "employeeforeignlanguage"
}
