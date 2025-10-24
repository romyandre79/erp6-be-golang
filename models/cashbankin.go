package models

import (
	"time"
)

type Cashbankin struct {
	Cashbankinid   int         `gorm:"column:cashbankinid;primaryKey" json:"cashbankinid"`
	Cashbankinno   string      `gorm:"column:cashbankinno" json:"cashbankinno"`
	Plantid        int         `gorm:"column:plantid" json:"plantid"`
	Cashbankindate time.Time   `gorm:"column:cashbankindate" json:"cashbankindate"`
	Notagirid      int         `gorm:"column:notagirid" json:"notagirid"`
	Addressbookid  int         `gorm:"column:addressbookid" json:"addressbookid"`
	Accountid      int         `gorm:"column:accountid" json:"accountid"`
	Headernote     string      `gorm:"column:headernote" json:"headernote"`
	Recordstatus   int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string      `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook    Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Account        Account     `gorm:"foreignKey:accountid;references:accountid"`
	Plant          Plant       `gorm:"foreignKey:plantid;references:plantid"`
}

func (Cashbankin) TableName() string {
	return "cashbankin"
}
