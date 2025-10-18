package models

import (
  "time"
)

type Invoicearsj struct {
	Invoicearsjid int `gorm:"column:invoicearsjid;primaryKey" json:"invoicearsjid"`
	Invoicearid int `gorm:"column:invoicearid" json:"invoicearid"`
	Giheaderid int `gorm:"column:giheaderid" json:"giheaderid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Giheader Giheader `gorm:"foreignKey:giheaderid;references:giheaderid"`
}

func (Invoicearsj) TableName() string {
	return "invoicearsj"
}
