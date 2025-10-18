package models

import (
  "time"
)

type Bookauthors struct {
	Bookauthorsid int `gorm:"column:bookauthorsid;primaryKey" json:"bookauthorsid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Authorname string `gorm:"column:authorname" json:"authorname"`
	Mobilephone *string `gorm:"column:mobilephone" json:"mobilephone"`
	Bookauthortypeid int `gorm:"column:bookauthortypeid" json:"bookauthortypeid"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Bookauthortype Bookauthortype `gorm:"foreignKey:bookauthortypeid;references:bookauthortypeid"`
}

func (Bookauthors) TableName() string {
	return "bookauthors"
}
