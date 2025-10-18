package models

import (

)

type Materisiswa struct {
	Materisiswaid int `gorm:"column:materisiswaid;primaryKey" json:"materisiswaid"`
	Materiid int `gorm:"column:materiid" json:"materiid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
}

func (Materisiswa) TableName() string {
	return "materisiswa"
}
