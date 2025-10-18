package models

import (
  "time"
)

type Employeeworkaccident struct {
	Employeeworkaccidentid int `gorm:"column:employeeworkaccidentid;primaryKey" json:"employeeworkaccidentid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Workaccidentdate time.Time `gorm:"column:workaccidentdate" json:"workaccidentdate"`
	Penanganan string `gorm:"column:penanganan" json:"penanganan"`
	Keteranganworkacc string `gorm:"column:keteranganworkacc" json:"keteranganworkacc"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeeworkaccident) TableName() string {
	return "employeeworkaccident"
}
