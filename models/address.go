package models

import (
	"time"
)

type Address struct {
	Addressid     int         `gorm:"column:addressid;primaryKey" json:"addressid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Addresstypeid int         `gorm:"column:addresstypeid" json:"addresstypeid"`
	Addressname   string      `gorm:"column:addressname" json:"addressname"`
	Cityid        int         `gorm:"column:cityid" json:"cityid"`
	Phoneno       string      `gorm:"column:phoneno" json:"phoneno"`
	Faxno         string      `gorm:"column:faxno" json:"faxno"`
	Lat           float64     `gorm:"column:lat" json:"lat"`
	Lng           float64     `gorm:"column:lng" json:"lng"`
	Foto          string      `gorm:"column:foto" json:"foto"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	City          City        `gorm:"foreignKey:cityid;references:cityid"`
	Addresstype   Addresstype `gorm:"foreignKey:addresstypeid;references:addresstypeid"`
}

func (Address) TableName() string {
	return "address"
}
