package models

import "time"

type Currency struct {
	Currencyid int `gorm:"column:currencyid;primaryKey" json:"currencyid"`
	Countryid int `gorm:"column:countryid" json:"countryid"`
	Currencyname string `gorm:"column:currencyname" json:"currencyname"`
	Symbol string `gorm:"column:symbol" json:"symbol"`
	I18n string `gorm:"column:i18n" json:"i18n"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Currency) TableName() string {
	return "currency"
}
