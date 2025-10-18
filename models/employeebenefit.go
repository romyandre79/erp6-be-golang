package models

import (
  "time"
)

type Employeebenefit struct {
	Employeebenefitid int `gorm:"column:employeebenefitid;primaryKey" json:"employeebenefitid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeebenefit) TableName() string {
	return "employeebenefit"
}
