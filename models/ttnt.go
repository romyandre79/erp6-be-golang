package models

import (
	"time"
)

type Ttnt struct {
	Ttntid        int         `gorm:"column:ttntid;primaryKey" json:"ttntid"`
	Ttntdate      time.Time   `gorm:"column:ttntdate" json:"ttntdate"`
	Ttntno        string      `gorm:"column:ttntno" json:"ttntno"`
	Plantid       int         `gorm:"column:plantid" json:"plantid"`
	Plantcode     string      `gorm:"column:plantcode" json:"plantcode"`
	Companyid     int         `gorm:"column:companyid" json:"companyid"`
	Companycode   string      `gorm:"column:companycode" json:"companycode"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Fullname      string      `gorm:"column:fullname" json:"fullname"`
	Description   string      `gorm:"column:description" json:"description"`
	Recordstatus  int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname    string      `gorm:"column:statusname" json:"statusname"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant         Plant       `gorm:"foreignKey:plantid;references:plantid"`
}

func (Ttnt) TableName() string {
	return "ttnt"
}
