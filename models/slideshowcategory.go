package models

import (
  "time"
)

type Slideshowcategory struct {
	Slideshowcategoryid int `gorm:"column:slideshowcategoryid;primaryKey" json:"slideshowcategoryid"`
	Slideshowid int `gorm:"column:slideshowid" json:"slideshowid"`
	Categoryid int `gorm:"column:categoryid" json:"categoryid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Category Category `gorm:"foreignKey:categoryid;references:categoryid"`
}

func (Slideshowcategory) TableName() string {
	return "slideshowcategory"
}
