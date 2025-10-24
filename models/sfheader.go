package models

import (
	"time"
)

type Sfheader struct {
	Sfheaderid      int         `gorm:"column:sfheaderid;primaryKey" json:"sfheaderid"`
	Sfno            string      `gorm:"column:sfno" json:"sfno"`
	Sfdate          time.Time   `gorm:"column:sfdate" json:"sfdate"`
	Plantid         int         `gorm:"column:plantid" json:"plantid"`
	Plantcode       string      `gorm:"column:plantcode" json:"plantcode"`
	Addressbookid   int         `gorm:"column:addressbookid" json:"addressbookid"`
	Pocustno        string      `gorm:"column:pocustno" json:"pocustno"`
	Addresstoid     int         `gorm:"column:addresstoid" json:"addresstoid"`
	Addresspayid    int         `gorm:"column:addresspayid" json:"addresspayid"`
	Salesid         int         `gorm:"column:salesid" json:"salesid"`
	Paymentmethodid int         `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Currencyid      int         `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate    float64     `gorm:"column:currencyrate" json:"currencyrate"`
	Isexport        int8        `gorm:"column:isexport" json:"isexport"`
	Issample        int8        `gorm:"column:issample" json:"issample"`
	Isavalan        int8        `gorm:"column:isavalan" json:"isavalan"`
	Headernote      string      `gorm:"column:headernote" json:"headernote"`
	Totalbefdisc    float64     `gorm:"column:totalbefdisc" json:"totalbefdisc"`
	Totalaftdisc    float64     `gorm:"column:totalaftdisc" json:"totalaftdisc"`
	Pendinganso     float64     `gorm:"column:pendinganso" json:"pendinganso"`
	Recordstatus    int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string      `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook     Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Addresspay      Address     `gorm:"foreignKey:addresspayid;references:addressid"`
	Addressto       Address     `gorm:"foreignKey:addresstoid;references:addressid"`
	Currency        Currency    `gorm:"foreignKey:currencyid;references:currencyid"`
	Plant           Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Sales           Sales       `gorm:"foreignKey:salesid;references:salesid"`
}

func (Sfheader) TableName() string {
	return "sfheader"
}
