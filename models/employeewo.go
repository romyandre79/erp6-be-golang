package models

import (
	"time"
)

type Employeewo struct {
	Employeewoid  int       `gorm:"column:employeewoid;primaryKey" json:"employeewoid"`
	Addressbookid int       `gorm:"column:addressbookid" json:"addressbookid"`
	Description   string    `gorm:"column:description" json:"description"`
	Company       string    `gorm:"column:company" json:"company"`
	Period        string    `gorm:"column:period" json:"period"`
	Isparklaring  int8      `gorm:"column:isparklaring" json:"isparklaring"`
	Recordstatus  int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeewo) TableName() string {
	return "employeewo"
}
