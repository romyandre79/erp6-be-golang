package models

import (
	"time"
)

type Cashbankjurnal struct {
	Cashbankjurnalid int       `gorm:"column:cashbankjurnalid;primaryKey" json:"cashbankjurnalid"`
	Cashbankid       int       `gorm:"column:cashbankid" json:"cashbankid"`
	Accountid        int       `gorm:"column:accountid" json:"accountid"`
	Accountname      string    `gorm:"column:accountname" json:"accountname"`
	Debit            float64   `gorm:"column:debit" json:"debit"`
	Credit           float64   `gorm:"column:credit" json:"credit"`
	Currencyid       int       `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue        float64   `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote       string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account          Account   `gorm:"foreignKey:accountid;references:accountid"`
	Cashbank         Cashbank  `gorm:"foreignKey:cashbankid;references:cashbankid"`
	Currency         Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Cashbankjurnal) TableName() string {
	return "cashbankjurnal"
}
