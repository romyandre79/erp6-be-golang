package models

import (
	"time"
)

type Packinglist struct {
	Packinglistid   int         `gorm:"column:packinglistid;primaryKey" json:"packinglistid"`
	Packinglistdate time.Time   `gorm:"column:packinglistdate" json:"packinglistdate"`
	Packinglistno   *string     `gorm:"column:packinglistno" json:"packinglistno"`
	Plantid         int         `gorm:"column:plantid" json:"plantid"`
	Plantcode       *string     `gorm:"column:plantcode" json:"plantcode"`
	Giheaderid      int         `gorm:"column:giheaderid" json:"giheaderid"`
	Gino            *string     `gorm:"column:gino" json:"gino"`
	Addressbookid   int         `gorm:"column:addressbookid" json:"addressbookid"`
	Fullname        *string     `gorm:"column:fullname" json:"fullname"`
	Supplierid      int         `gorm:"column:supplierid" json:"supplierid"`
	Suppliername    *string     `gorm:"column:suppliername" json:"suppliername"`
	Addresstoid     int         `gorm:"column:addresstoid" json:"addresstoid"`
	Addressname     *string     `gorm:"column:addressname" json:"addressname"`
	Soheaderid      int         `gorm:"column:soheaderid" json:"soheaderid"`
	Sono            *string     `gorm:"column:sono" json:"sono"`
	Nomobil         *string     `gorm:"column:nomobil" json:"nomobil"`
	Pebno           *string     `gorm:"column:pebno" json:"pebno"`
	Sopir           *string     `gorm:"column:sopir" json:"sopir"`
	Pocustno        *string     `gorm:"column:pocustno" json:"pocustno"`
	Headernote      *string     `gorm:"column:headernote" json:"headernote"`
	Recordstatus    int         `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string      `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook     Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Address         Address     `gorm:"foreignKey:addresstoid;references:addressid"`
	Plant           Plant       `gorm:"foreignKey:plantid;references:plantid"`
	Soheader        Soheader    `gorm:"foreignKey:soheaderid;references:soheaderid"`
	Supplier        Addressbook `gorm:"foreignKey:supplierid;references:addressbookid"`
}

func (Packinglist) TableName() string {
	return "packinglist"
}
