package models

import (
  "time"
)

type Employeekontrak struct {
	Employeekontrakid int `gorm:"column:employeekontrakid;primaryKey" json:"employeekontrakid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeekontrak) TableName() string {
	return "employeekontrak"
}
