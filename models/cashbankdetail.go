package models

import (
	"time"
)

type Cashbankdetail struct {
	Cashbankdetailid int       `gorm:"column:cashbankdetailid;primaryKey" json:"cashbankdetailid"`
	Cashbankid       int       `gorm:"column:cashbankid" json:"cashbankid"`
	Slocid           int       `gorm:"column:slocid" json:"slocid"`
	Accountid        int       `gorm:"column:accountid" json:"accountid"`
	Amount           float64   `gorm:"column:amount" json:"amount"`
	Currencyid       int       `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue        float64   `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote       string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account          Account   `gorm:"foreignKey:accountid;references:accountid"`
	Currency         Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
	Sloc             Sloc      `gorm:"foreignKey:slocid;references:slocid"`
}

func (Cashbankdetail) TableName() string {
	return "cashbankdetail"
}
