package models

import (

)

type Bookwish struct {
	Bookwishid int `gorm:"column:bookwishid;primaryKey" json:"bookwishid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Book Book `gorm:"foreignKey:bookid;references:bookid"`
}

func (Bookwish) TableName() string {
	return "bookwish"
}
