package models

import (
  "time"
)

type Employeeoverdet struct {
	Employeeoverdetid int `gorm:"column:employeeoverdetid;primaryKey" json:"employeeoverdetid"`
	Employeeoverid int `gorm:"column:employeeoverid" json:"employeeoverid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Overtimestart time.Time `gorm:"column:overtimestart" json:"overtimestart"`
	Overtimeend time.Time `gorm:"column:overtimeend" json:"overtimeend"`
	Reason *string `gorm:"column:reason" json:"reason"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeeoverdet) TableName() string {
	return "employeeoverdet"
}
