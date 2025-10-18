package models

import (
  "time"
)

type Employeewstatus struct {
	Employeewstatusid int `gorm:"column:employeewstatusid;primaryKey" json:"employeewstatusid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Employeestatusid int `gorm:"column:employeestatusid" json:"employeestatusid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Employeestatus Employeestatus `gorm:"foreignKey:employeestatusid;references:employeestatusid"`
}

func (Employeewstatus) TableName() string {
	return "employeewstatus"
}
