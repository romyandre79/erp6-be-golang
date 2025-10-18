package models

import (
  "time"
)

type Notagrr struct {
	Notagrrid int `gorm:"column:notagrrid;primaryKey" json:"notagrrid"`
	Notagrrno *string `gorm:"column:notagrrno" json:"notagrrno"`
	Notagrrdate *time.Time `gorm:"column:notagrrdate" json:"notagrrdate"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Invoiceaptaxno string `gorm:"column:invoiceaptaxno" json:"invoiceaptaxno"`
	Poheaderid int `gorm:"column:poheaderid" json:"poheaderid"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Poheader Poheader `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Notagrr) TableName() string {
	return "notagrr"
}
