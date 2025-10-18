package models

import (
  "time"
)

type Fajurnal struct {
	Fajurnalid int `gorm:"column:fajurnalid;primaryKey" json:"fajurnalid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Fixassetid int `gorm:"column:fixassetid" json:"fixassetid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Debet float64 `gorm:"column:debet" json:"debet"`
	Credit float64 `gorm:"column:credit" json:"credit"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Susutdate time.Time `gorm:"column:susutdate" json:"susutdate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Fajurnal) TableName() string {
	return "fajurnal"
}
