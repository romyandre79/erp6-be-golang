package models

import "time"

type Province struct {
	Provinceid int `gorm:"column:provinceid;primaryKey" json:"provinceid"`
	Countryid int `gorm:"column:countryid" json:"countryid"`
	Provincecode string `gorm:"column:provincecode" json:"provincecode"`
	Provincename string `gorm:"column:provincename" json:"provincename"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Country *Country `gorm:"foreignKey:countryid;references:countryid"`
}

func (Province) TableName() string {
	return "province"
}
