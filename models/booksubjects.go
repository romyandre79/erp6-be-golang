package models

import (

)

type Booksubjects struct {
	Booksubjectsid int `gorm:"column:booksubjectsid;primaryKey" json:"booksubjectsid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Booksubjectid int `gorm:"column:booksubjectid" json:"booksubjectid"`
	Book Book `gorm:"foreignKey:bookid;references:bookid"`
	Booksubject Booksubject `gorm:"foreignKey:booksubjectid;references:booksubjectid"`
}

func (Booksubjects) TableName() string {
	return "booksubjects"
}
