package models

import (
  "time"
)

type Photoscategory struct {
	Photoscategoryid int `gorm:"column:photoscategoryid;primaryKey" json:"photoscategoryid"`
	Photosid int `gorm:"column:photosid" json:"photosid"`
	Categoryid int `gorm:"column:categoryid" json:"categoryid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Category Category `gorm:"foreignKey:categoryid;references:categoryid"`
	Photos Photos `gorm:"foreignKey:photosid;references:photosid"`
}

func (Photoscategory) TableName() string {
	return "photoscategory"
}
