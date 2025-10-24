package models

import (
  "time"
)

type Expeditionarjurnal struct {
	Expeditionarjurnalid int `gorm:"column:expeditionarjurnalid;primaryKey" json:"expeditionarjurnalid"`
	Expeditionarid int `gorm:"column:expeditionarid" json:"expeditionarid"`
	Cashbankid int `gorm:"column:cashbankid" json:"cashbankid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote string `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Cashbank Cashbank `gorm:"foreignKey:cashbankid;references:cashbankid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Expeditionarjurnal) TableName() string {
	return "expeditionarjurnal"
}
