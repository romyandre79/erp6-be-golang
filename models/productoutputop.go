package models

import (
  "time"
)

type Productoutputop struct {
	Productoutputopid int `gorm:"column:productoutputopid;primaryKey" json:"productoutputopid"`
	Productoutputfgid int `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Employeeid int `gorm:"column:employeeid" json:"employeeid"`
	Employeename string `gorm:"column:employeename" json:"employeename"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee Employee `gorm:"foreignKey:employeeid;references:employeeid"`
	Productoutputfg Productoutputfg `gorm:"foreignKey:productoutputfgid;references:productoutputfgid"`
}

func (Productoutputop) TableName() string {
	return "productoutputop"
}
