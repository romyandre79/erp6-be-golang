package models

import (
	"time"
)

type Grheader struct {
	Grheaderid    int         `gorm:"column:grheaderid;primaryKey" json:"grheaderid"`
	Grdate        time.Time   `gorm:"column:grdate" json:"grdate"`
	Grno          string      `gorm:"column:grno" json:"grno"`
	Isjasa        int8        `gorm:"column:isjasa" json:"isjasa"`
	Plantid       int         `gorm:"column:plantid" json:"plantid"`
	Poheaderid    int         `gorm:"column:poheaderid" json:"poheaderid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Sjsupplier    string      `gorm:"column:sjsupplier" json:"sjsupplier"`
	Pibno         string      `gorm:"column:pibno" json:"pibno"`
	Kendaraanno   string      `gorm:"column:kendaraanno" json:"kendaraanno"`
	Supir         string      `gorm:"column:supir" json:"supir"`
	Headernote    string      `gorm:"column:headernote" json:"headernote"`
	Isinvap       int8        `gorm:"column:isinvap" json:"isinvap"`
	Recordstatus  int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname    string      `gorm:"column:statusname" json:"statusname"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant         Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Poheader      Poheader    `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Grheader) TableName() string {
	return "grheader"
}
