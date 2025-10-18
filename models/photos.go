package models

import (
  "time"
)

type Photos struct {
	Photosid int `gorm:"column:photosid;primaryKey" json:"photosid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Photostitle string `gorm:"column:photostitle" json:"photostitle"`
	Photossubtitle string `gorm:"column:photossubtitle" json:"photossubtitle"`
	Photosdesc string `gorm:"column:photosdesc" json:"photosdesc"`
	Photosurl string `gorm:"column:photosurl" json:"photosurl"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Photos) TableName() string {
	return "photos"
}
