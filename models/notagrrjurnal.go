package models

import (
  "time"
)

type Notagrrjurnal struct {
	Notagrrjurnalid int `gorm:"column:notagrrjurnalid;primaryKey" json:"notagrrjurnalid"`
	Notagrrid int `gorm:"column:notagrrid" json:"notagrrid"`
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

func (Notagrrjurnal) TableName() string {
	return "notagrrjurnal"
}
