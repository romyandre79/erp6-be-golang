package models

import (
  "time"
)

type Abstrans struct {
	Abstransid int `gorm:"column:abstransid;primaryKey" json:"abstransid"`
	Employeeid int `gorm:"column:employeeid" json:"employeeid"`
	Clockdatetime time.Time `gorm:"column:clockdatetime" json:"clockdatetime"`
	Description *string `gorm:"column:description" json:"description"`
	Status *string `gorm:"column:status" json:"status"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Employee Employee `gorm:"foreignKey:employeeid;references:employeeid"`
}

func (Abstrans) TableName() string {
	return "abstrans"
}
