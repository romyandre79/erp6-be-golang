package models

import (
	"time"
)

type Invoicearpl struct {
	Invoicearplid int         `gorm:"column:invoicearplid;primaryKey" json:"invoicearplid"`
	Invoicearid   int         `gorm:"column:invoicearid" json:"invoicearid"`
	Invoicearsjid int         `gorm:"column:invoicearsjid" json:"invoicearsjid"`
	Packinglistid int         `gorm:"column:packinglistid" json:"packinglistid"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Invoicearsj   Invoicearsj `gorm:"foreignKey:invoicearsjid;references:invoicearsjid"`
	Packinglist   Packinglist `gorm:"foreignKey:packinglistid;references:packinglistid"`
}

func (Invoicearpl) TableName() string {
	return "invoicearpl"
}
