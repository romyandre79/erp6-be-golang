package models

import (
	"time"
)

type Certoa struct {
	Certoaid       int            `gorm:"column:certoaid;primaryKey" json:"certoaid"`
	Certoano       string         `gorm:"column:certoano" json:"certoano"`
	Certoadate     time.Time      `gorm:"column:certoadate" json:"certoadate"`
	Plantid        int            `gorm:"column:plantid" json:"plantid"`
	Productiondate time.Time      `gorm:"column:productiondate" json:"productiondate"`
	Soheaderid     int            `gorm:"column:soheaderid" json:"soheaderid"`
	Addressbookid  int            `gorm:"column:addressbookid" json:"addressbookid"`
	Productid      int            `gorm:"column:productid" json:"productid"`
	Bomid          int            `gorm:"column:bomid" json:"bomid"`
	Jumkirim       float64        `gorm:"column:jumkirim" json:"jumkirim"`
	Status         int8           `gorm:"column:status" json:"status"`
	Description    string         `gorm:"column:description" json:"description"`
	Recordstatus   int8           `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string         `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Addressbook    Addressbook    `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Billofmaterial Billofmaterial `gorm:"foreignKey:bomid;references:bomid"`
	Plant          Plant          `gorm:"foreignKey:plantid;references:plantid"`
	Product        Product        `gorm:"foreignKey:productid;references:productid"`
	Soheader       Soheader       `gorm:"foreignKey:soheaderid;references:soheaderid"`
}

func (Certoa) TableName() string {
	return "certoa"
}
