package models

import (
  "time"
)

type Sales struct {
	Salesid int `gorm:"column:salesid;primaryKey" json:"salesid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Limitsampleamount float64 `gorm:"column:limitsampleamount" json:"limitsampleamount"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Sales) TableName() string {
	return "sales"
}
