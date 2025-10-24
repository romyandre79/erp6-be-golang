package models

import (
  "time"
)

type Invoiceapgr struct {
	Invoiceapgrid int `gorm:"column:invoiceapgrid;primaryKey" json:"invoiceapgrid"`
	Invoiceapid int `gorm:"column:invoiceapid" json:"invoiceapid"`
	Grheaderid int `gorm:"column:grheaderid" json:"grheaderid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Grheader Grheader `gorm:"foreignKey:grheaderid;references:grheaderid"`
}

func (Invoiceapgr) TableName() string {
	return "invoiceapgr"
}
