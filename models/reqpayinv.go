package models

import (
	"time"
)

type Reqpayinv struct {
	Reqpayinvid   int         `gorm:"column:reqpayinvid;primaryKey" json:"reqpayinvid"`
	Reqpayid      int         `gorm:"column:reqpayid" json:"reqpayid"`
	Invoiceapid   int         `gorm:"column:invoiceapid" json:"invoiceapid"`
	Notagrrid     int         `gorm:"column:notagrrid" json:"notagrrid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Duedate       time.Time   `gorm:"column:duedate" json:"duedate"`
	Amount        float64     `gorm:"column:amount" json:"amount"`
	Currencyid    int         `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue     float64     `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote    string      `gorm:"column:detailnote" json:"detailnote"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Currency      Currency    `gorm:"foreignKey:currencyid;references:currencyid"`
	Invoiceap     Invoiceap   `gorm:"foreignKey:invoiceapid;references:invoiceapid"`
	Notagrr       Notagrr     `gorm:"foreignKey:notagrrid;references:notagrrid"`
}

func (Reqpayinv) TableName() string {
	return "reqpayinv"
}
