package models

import (
  "time"
)

type Invoicear struct {
	Invoicearid int `gorm:"column:invoicearid;primaryKey" json:"invoicearid"`
	Invoiceardate time.Time `gorm:"column:invoiceardate" json:"invoiceardate"`
	Invoicearno *string `gorm:"column:invoicearno" json:"invoicearno"`
	Doctype int8 `gorm:"column:doctype" json:"doctype"`
	Soheaderid *int `gorm:"column:soheaderid" json:"soheaderid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Companycode string `gorm:"column:companycode" json:"companycode"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Contractno string `gorm:"column:contractno" json:"contractno"`
	Invoiceartaxno string `gorm:"column:invoiceartaxno" json:"invoiceartaxno"`
	Taxno string `gorm:"column:taxno" json:"taxno"`
	Addressname string `gorm:"column:addressname" json:"addressname"`
	Paymentmethodid int `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Duedate time.Time `gorm:"column:duedate" json:"duedate"`
	Cashbankid *int `gorm:"column:cashbankid" json:"cashbankid"`
	Dpamount float64 `gorm:"column:dpamount" json:"dpamount"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Iscbin int8 `gorm:"column:iscbin" json:"iscbin"`
	Total *float64 `gorm:"column:total" json:"total"`
	Cashbankamount *float64 `gorm:"column:cashbankamount" json:"cashbankamount"`
	Repeatby *int8 `gorm:"column:repeatby" json:"repeatby"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Paymentmethod Paymentmethod `gorm:"foreignKey:paymentmethodid;references:paymentmethodid"`
	Soheader Soheader `gorm:"foreignKey:soheaderid;references:soheaderid"`
}

func (Invoicear) TableName() string {
	return "invoicear"
}
