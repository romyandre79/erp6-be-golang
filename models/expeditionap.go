package models

import (
	"time"
)

type Expeditionap struct {
	Expeditionapid   int         `gorm:"column:expeditionapid;primaryKey" json:"expeditionapid"`
	Expeditionapno   *string     `gorm:"column:expeditionapno" json:"expeditionapno"`
	Plantid          int         `gorm:"column:plantid" json:"plantid"`
	Expeditionapdate time.Time   `gorm:"column:expeditionapdate" json:"expeditionapdate"`
	Poheaderid       int         `gorm:"column:poheaderid" json:"poheaderid"`
	Addressbookid    int         `gorm:"column:addressbookid" json:"addressbookid"`
	Addressbookexpid int         `gorm:"column:addressbookexpid" json:"addressbookexpid"`
	Amount           float64     `gorm:"column:amount" json:"amount"`
	Headernote       string      `gorm:"column:headernote" json:"headernote"`
	Recordstatus     int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname       string      `gorm:"column:statusname" json:"statusname"`
	Updatedate       time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook      Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Addressbookexp   Addressbook `gorm:"foreignKey:addressbookexpid;references:addressbookid"`
	Plant            Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Poheader         Poheader    `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Expeditionap) TableName() string {
	return "expeditionap"
}
