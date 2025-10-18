package models

import (
  "time"
)

type Expeditionapjurnal struct {
	Expeditionapjurnalid int `gorm:"column:expeditionapjurnalid;primaryKey" json:"expeditionapjurnalid"`
	Expeditionapid int `gorm:"column:expeditionapid" json:"expeditionapid"`
	Cashbankid int `gorm:"column:cashbankid" json:"cashbankid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote string `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Expeditionapjurnal) TableName() string {
	return "expeditionapjurnal"
}
