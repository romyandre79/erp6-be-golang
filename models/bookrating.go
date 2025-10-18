package models

import (

)

type Bookrating struct {
	Bookratingid int `gorm:"column:bookratingid;primaryKey" json:"bookratingid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Bookmemberid int `gorm:"column:bookmemberid" json:"bookmemberid"`
	Book Book `gorm:"foreignKey:bookid;references:bookid"`
	Bookmember Bookmember `gorm:"foreignKey:bookmemberid;references:bookmemberid"`
}

func (Bookrating) TableName() string {
	return "bookrating"
}
