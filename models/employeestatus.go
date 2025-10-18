package models

import (
  "time"
)

type Employeestatus struct {
	Employeestatusid int `gorm:"column:employeestatusid;primaryKey" json:"employeestatusid"`
	Employeestatusname string `gorm:"column:employeestatusname" json:"employeestatusname"`
	Description string `gorm:"column:description" json:"description"`
	Taxvalue float64 `gorm:"column:taxvalue" json:"taxvalue"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeestatus) TableName() string {
	return "employeestatus"
}
