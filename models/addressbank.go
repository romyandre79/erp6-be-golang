package models

import (

)

type Addressbank struct {
	Addressbookbankid int `gorm:"column:addressbookbankid;primaryKey" json:"addressbookbankid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Bankid int `gorm:"column:bankid" json:"bankid"`
	Bankaccountno string `gorm:"column:bankaccountno" json:"bankaccountno"`
	Bankownername string `gorm:"column:bankownername" json:"bankownername"`
	Bank Bank `gorm:"foreignKey:bankid;references:bankid"`
}

func (Addressbank) TableName() string {
	return "addressbank"
}
