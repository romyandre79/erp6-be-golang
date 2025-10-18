package models

import (
  "time"
)

type Slideshow struct {
	Slideshowid int `gorm:"column:slideshowid;primaryKey" json:"slideshowid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Slidepic string `gorm:"column:slidepic" json:"slidepic"`
	Slidetitle string `gorm:"column:slidetitle" json:"slidetitle"`
	Slidedesc string `gorm:"column:slidedesc" json:"slidedesc"`
	Slideurl string `gorm:"column:slideurl" json:"slideurl"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Slideshow) TableName() string {
	return "slideshow"
}
