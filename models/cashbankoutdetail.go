package models

import (
	"time"
)

type Cashbankoutdetail struct {
	Cashbankoutdetailid int       `gorm:"column:cashbankoutdetailid;primaryKey" json:"cashbankoutdetailid"`
	Cashbankoutid       int       `gorm:"column:cashbankoutid" json:"cashbankoutid"`
	Addressbookid       int       `gorm:"column:addressbookid" json:"addressbookid"`
	Invoiceapid         int       `gorm:"column:invoiceapid" json:"invoiceapid"`
	Notagrrid           int       `gorm:"column:notagrrid" json:"notagrrid"`
	Amount              float64   `gorm:"column:amount" json:"amount"`
	Currencyid          int       `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue           float64   `gorm:"column:ratevalue" json:"ratevalue"`
	Chequeid            int       `gorm:"column:chequeid" json:"chequeid"`
	Nobuktipotong       string    `gorm:"column:nobuktipotong" json:"nobuktipotong"`
	Detailnote          string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate          time.Time `gorm:"column:updatedate" json:"updatedate"`
	Invoiceap           Invoiceap `gorm:"foreignKey:invoiceapid;references:invoiceapid"`
	Notagrr             Notagrr   `gorm:"foreignKey:notagrrid;references:notagrrid"`
}

func (Cashbankoutdetail) TableName() string {
	return "cashbankoutdetail"
}
