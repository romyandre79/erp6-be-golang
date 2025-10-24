package models

import (
  "time"
)

type Qctest struct {
	Qctestid int `gorm:"column:qctestid;primaryKey" json:"qctestid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
}

func (Qctest) TableName() string {
	return "qctest"
}
