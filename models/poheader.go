package models

import (
  "time"
)

type Poheader struct {
	Poheaderid int `gorm:"column:poheaderid;primaryKey" json:"poheaderid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Pono *string `gorm:"column:pono" json:"pono"`
	Podate time.Time `gorm:"column:podate" json:"podate"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Addresscontactid *int `gorm:"column:addresscontactid" json:"addresscontactid"`
	Paymentmethodid int `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Isimport int8 `gorm:"column:isimport" json:"isimport"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Printke *float64 `gorm:"column:printke" json:"printke"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Isjasa *int8 `gorm:"column:isjasa" json:"isjasa"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Paymentmethod Paymentmethod `gorm:"foreignKey:paymentmethodid;references:paymentmethodid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Poheader) TableName() string {
	return "poheader"
}
