package models

import (
  "time"
)

type Productplan struct {
	Productplanid int `gorm:"column:productplanid;primaryKey" json:"productplanid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Soheaderid *int `gorm:"column:soheaderid" json:"soheaderid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Parentplanid *int `gorm:"column:parentplanid" json:"parentplanid"`
	Productplanno *string `gorm:"column:productplanno" json:"productplanno"`
	Productplandate time.Time `gorm:"column:productplandate" json:"productplandate"`
	Description *string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Productplan) TableName() string {
	return "productplan"
}
