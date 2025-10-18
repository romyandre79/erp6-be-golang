package models

import (
  "time"
)

type Cashbankinacc struct {
	Cashbankinaccid int `gorm:"column:cashbankinaccid;primaryKey" json:"cashbankinaccid"`
	Cashbankinid *int `gorm:"column:cashbankinid" json:"cashbankinid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Debit float64 `gorm:"column:debit" json:"debit"`
	Credit float64 `gorm:"column:credit" json:"credit"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Itemnote string `gorm:"column:itemnote" json:"itemnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Cashbankinacc) TableName() string {
	return "cashbankinacc"
}
