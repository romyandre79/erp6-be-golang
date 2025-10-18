package models

import (
  "time"
)

type Addresscompany struct {
	Addresscompanyid int `gorm:"column:addresscompanyid;primaryKey" json:"addresscompanyid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Addresscompany) TableName() string {
	return "addresscompany"
}
