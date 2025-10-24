package models

import (
	"time"
)

type Mapelkelasdiskusi struct {
	Mapelkelasdiskusiid int                `gorm:"column:mapelkelasdiskusiid;primaryKey" json:"mapelkelasdiskusiid"`
	Mapelkelasid        int                `gorm:"column:mapelkelasid" json:"mapelkelasid"`
	Addressbookid       int                `gorm:"column:addressbookid" json:"addressbookid"`
	Content             string             `gorm:"column:content" json:"content"`
	Replyid             int                `gorm:"column:replyid" json:"replyid"`
	Chatdate            time.Time          `gorm:"column:chatdate" json:"chatdate"`
	Addressbook         Addressbook        `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Mapelkelas          Mapelkelas         `gorm:"foreignKey:mapelkelasid;references:mapelkelasid"`
	Reply               *Mapelkelasdiskusi `gorm:"foreignKey:replyid;references:mapelkelasdiskusiid"`
}

func (Mapelkelasdiskusi) TableName() string {
	return "mapelkelasdiskusi"
}
