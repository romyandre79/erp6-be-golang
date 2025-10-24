package models

import (
	"time"
)

type Country struct {
	Countryid    int       `gorm:"column:countryid;primaryKey" json:"countryid"`
	Countrycode  string    `gorm:"column:countrycode" json:"countrycode"`
	Countryname  string    `gorm:"column:countryname" json:"countryname"`
	Recordstatus int       `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Country) TableName() string {
	return "country"
}
