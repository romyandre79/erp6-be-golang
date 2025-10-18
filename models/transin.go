package models

import (
  "time"
)

type Transin struct {
	Transinid int `gorm:"column:transinid;primaryKey" json:"transinid"`
	Docdate time.Time `gorm:"column:docdate" json:"docdate"`
	Transinno *string `gorm:"column:transinno" json:"transinno"`
	Employeeid int `gorm:"column:employeeid" json:"employeeid"`
	Permitinid int `gorm:"column:permitinid" json:"permitinid"`
	Description *string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee Employee `gorm:"foreignKey:employeeid;references:employeeid"`
	Permitin Permitin `gorm:"foreignKey:permitinid;references:permitinid"`
}

func (Transin) TableName() string {
	return "transin"
}
