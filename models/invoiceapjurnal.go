package models

import (
  "time"
)

type Invoiceapjurnal struct {
	Invoiceapjurnalid int `gorm:"column:invoiceapjurnalid;primaryKey" json:"invoiceapjurnalid"`
	Invoiceapid int `gorm:"column:invoiceapid" json:"invoiceapid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Accountname string `gorm:"column:accountname" json:"accountname"`
	Debit float64 `gorm:"column:debit" json:"debit"`
	Credit float64 `gorm:"column:credit" json:"credit"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote string `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Invoiceapjurnal) TableName() string {
	return "invoiceapjurnal"
}
