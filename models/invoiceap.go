package models

import (
	"time"
)

type Invoiceap struct {
	Invoiceapid     int           `gorm:"column:invoiceapid;primaryKey" json:"invoiceapid"`
	Doctype         int           `gorm:"column:doctype" json:"doctype"`
	Invoiceapdate   time.Time     `gorm:"column:invoiceapdate" json:"invoiceapdate"`
	Invoiceapno     string        `gorm:"column:invoiceapno" json:"invoiceapno"`
	Poheaderid      int           `gorm:"column:poheaderid" json:"poheaderid"`
	Pono            string        `gorm:"column:pono" json:"pono"`
	Companyid       int           `gorm:"column:companyid" json:"companyid"`
	Companycode     string        `gorm:"column:companycode" json:"companycode"`
	Plantid         int           `gorm:"column:plantid" json:"plantid"`
	Plantcode       string        `gorm:"column:plantcode" json:"plantcode"`
	Addressbookid   int           `gorm:"column:addressbookid" json:"addressbookid"`
	Contractno      string        `gorm:"column:contractno" json:"contractno"`
	Invoiceaptaxno  string        `gorm:"column:invoiceaptaxno" json:"invoiceaptaxno"`
	Taxno           string        `gorm:"column:taxno" json:"taxno"`
	Nobuktipotong   string        `gorm:"column:nobuktipotong" json:"nobuktipotong"`
	Receiptdate     time.Time     `gorm:"column:receiptdate" json:"receiptdate"`
	Beritaacara     string        `gorm:"column:beritaacara" json:"beritaacara"`
	Paymentmethodid int           `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Paycode         string        `gorm:"column:paycode" json:"paycode"`
	Duedate         time.Time     `gorm:"column:duedate" json:"duedate"`
	Discount        float64       `gorm:"column:discount" json:"discount"`
	Cashbankid      int           `gorm:"column:cashbankid" json:"cashbankid"`
	Dpamount        float64       `gorm:"column:dpamount" json:"dpamount"`
	Headernote      string        `gorm:"column:headernote" json:"headernote"`
	Isreqpay        int8          `gorm:"column:isreqpay" json:"isreqpay"`
	Payamount       float64       `gorm:"column:payamount" json:"payamount"`
	Amount          float64       `gorm:"column:amount" json:"amount"`
	Recordstatus    int8          `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string        `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Addressbook     Addressbook   `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Paymentmethod   Paymentmethod `gorm:"foreignKey:paymentmethodid;references:paymentmethodid"`
	Plant           Plant         `gorm:"foreignKey:plantid;references:plantid"`
	Poheader        Poheader      `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Invoiceap) TableName() string {
	return "invoiceap"
}
