package models

import (

)

type Companybank struct {
	Companybankid int `gorm:"column:companybankid;primaryKey" json:"companybankid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Bankid int `gorm:"column:bankid" json:"bankid"`
	Bankaccountno string `gorm:"column:bankaccountno" json:"bankaccountno"`
	Bankowner string `gorm:"column:bankowner" json:"bankowner"`
}

func (Companybank) TableName() string {
	return "companybank"
}
