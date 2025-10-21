package models

import (
	"time"
)

type Notagir struct {
	Notagirid      int         `gorm:"column:notagirid;primaryKey" json:"notagirid"`
	Notagirno      string      `gorm:"column:notagirno" json:"notagirno"`
	Notagirdate    time.Time   `gorm:"column:notagirdate" json:"notagirdate"`
	Plantid        int         `gorm:"column:plantid" json:"plantid"`
	Addressbookid  int         `gorm:"column:addressbookid" json:"addressbookid"`
	Soheaderid     int         `gorm:"column:soheaderid" json:"soheaderid"`
	Pocustno       string      `gorm:"column:pocustno" json:"pocustno"`
	Gireturid      int         `gorm:"column:gireturid" json:"gireturid"`
	Invoiceartaxno string      `gorm:"column:invoiceartaxno" json:"invoiceartaxno"`
	Headernote     string      `gorm:"column:headernote" json:"headernote"`
	Recordstatus   int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string      `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook    Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant          Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Soheader       Soheader    `gorm:"foreignKey:soheaderid;references:soheaderid"`
}

func (Notagir) TableName() string {
	return "notagir"
}
