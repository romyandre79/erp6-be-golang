package models

import (

)

type Pesertappdbdoc struct {
	Pesertappdbdocid int `gorm:"column:pesertappdbdocid;primaryKey" json:"pesertappdbdocid"`
	Pesertappdbid int `gorm:"column:pesertappdbid" json:"pesertappdbid"`
	Document string `gorm:"column:document" json:"document"`
	Filedoc string `gorm:"column:filedoc" json:"filedoc"`
}

func (Pesertappdbdoc) TableName() string {
	return "pesertappdbdoc"
}
