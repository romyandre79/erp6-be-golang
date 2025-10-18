package models

import (
  "time"
)

type Employeesp struct {
	Employeespid int `gorm:"column:employeespid;primaryKey" json:"employeespid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Spdate time.Time `gorm:"column:spdate" json:"spdate"`
	Expireddatesp int `gorm:"column:expireddatesp" json:"expireddatesp"`
	Keterangansp string `gorm:"column:keterangansp" json:"keterangansp"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeesp) TableName() string {
	return "employeesp"
}
