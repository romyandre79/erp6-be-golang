package models

import (
  "time"
)

type Notagirjurnal struct {
	Notagirjurnalid int `gorm:"column:notagirjurnalid;primaryKey" json:"notagirjurnalid"`
	Notagirid int `gorm:"column:notagirid" json:"notagirid"`
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

func (Notagirjurnal) TableName() string {
	return "notagirjurnal"
}
