package models

import (
  "time"
)

type Taxacc struct {
	Taxaccid int `gorm:"column:taxaccid;primaryKey" json:"taxaccid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Taxacc) TableName() string {
	return "taxacc"
}
