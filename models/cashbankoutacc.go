package models

import (
  "time"
)

type Cashbankoutacc struct {
	Cashbankoutaccid int `gorm:"column:cashbankoutaccid;primaryKey" json:"cashbankoutaccid"`
	Cashbankoutid *int `gorm:"column:cashbankoutid" json:"cashbankoutid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Debit int `gorm:"column:debit" json:"debit"`
	Credit int `gorm:"column:credit" json:"credit"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Itemnote string `gorm:"column:itemnote" json:"itemnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Cashbankoutacc) TableName() string {
	return "cashbankoutacc"
}
