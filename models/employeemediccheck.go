package models

import (
  "time"
)

type Employeemediccheck struct {
	Employeemediccheckid int `gorm:"column:employeemediccheckid;primaryKey" json:"employeemediccheckid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeemediccheck) TableName() string {
	return "employeemediccheck"
}
