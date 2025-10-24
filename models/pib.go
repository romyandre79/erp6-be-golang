package models

import (
  "time"
)

type Pib struct {
	Pibid int `gorm:"column:pibid;primaryKey" json:"pibid"`
	Poheaderid int `gorm:"column:poheaderid" json:"poheaderid"`
	Pibno string `gorm:"column:pibno" json:"pibno"`
	Lcno string `gorm:"column:lcno" json:"lcno"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Poheader Poheader `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Pib) TableName() string {
	return "pib"
}
