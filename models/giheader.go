package models

import (
	"time"
)

type Giheader struct {
	Giheaderid    int         `gorm:"column:giheaderid;primaryKey" json:"giheaderid"`
	Gidate        time.Time   `gorm:"column:gidate" json:"gidate"`
	Gino          *string     `gorm:"column:gino" json:"gino"`
	Plantid       int         `gorm:"column:plantid" json:"plantid"`
	Plantcode     string      `gorm:"column:plantcode" json:"plantcode"`
	Companyname   *string     `gorm:"column:companyname" json:"companyname"`
	Soheaderid    int         `gorm:"column:soheaderid" json:"soheaderid"`
	Pocustno      *string     `gorm:"column:pocustno" json:"pocustno"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Customername  string      `gorm:"column:customername" json:"customername"`
	Supplierid    int         `gorm:"column:supplierid" json:"supplierid"`
	Suppliername  string      `gorm:"column:suppliername" json:"suppliername"`
	Nomobil       string      `gorm:"column:nomobil" json:"nomobil"`
	Sopir         string      `gorm:"column:sopir" json:"sopir"`
	Addresstoid   int         `gorm:"column:addresstoid" json:"addresstoid"`
	Pebno         *string     `gorm:"column:pebno" json:"pebno"`
	Headernote    *string     `gorm:"column:headernote" json:"headernote"`
	Recordstatus  int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname    string      `gorm:"column:statusname" json:"statusname"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Address       Address     `gorm:"foreignKey:addresstoid;references:addressid"`
	Plant         Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Soheader      Soheader    `gorm:"foreignKey:soheaderid;references:soheaderid"`
	Supplier      Addressbook `gorm:"foreignKey:supplierid;references:addressbookid"`
}

func (Giheader) TableName() string {
	return "giheader"
}
