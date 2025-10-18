package models

import (
	"time"
)

type Expeditionar struct {
	Expeditionarid   int         `gorm:"column:expeditionarid;primaryKey" json:"expeditionarid"`
	Expeditionarno   *string     `gorm:"column:expeditionarno" json:"expeditionarno"`
	Plantid          int         `gorm:"column:plantid" json:"plantid"`
	Expeditionardate time.Time   `gorm:"column:expeditionardate" json:"expeditionardate"`
	Soheaderid       int         `gorm:"column:soheaderid" json:"soheaderid"`
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
	Soheader         Soheader    `gorm:"foreignKey:soheaderid;references:soheaderid"`
}

func (Expeditionar) TableName() string {
	return "expeditionar"
}
