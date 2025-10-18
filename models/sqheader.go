package models

import (
  "time"
)

type Sqheader struct {
	Sqheaderid int `gorm:"column:sqheaderid;primaryKey" json:"sqheaderid"`
	Sqno *string `gorm:"column:sqno" json:"sqno"`
	Sqdate time.Time `gorm:"column:sqdate" json:"sqdate"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Salesid *int `gorm:"column:salesid" json:"salesid"`
	Paymentmethodid int `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Paymentmethod Paymentmethod `gorm:"foreignKey:paymentmethodid;references:paymentmethodid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Sales Sales `gorm:"foreignKey:salesid;references:salesid"`
}

func (Sqheader) TableName() string {
	return "sqheader"
}
