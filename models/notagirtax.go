package models

import (
  "time"
)

type Notagirtax struct {
	Notagirtaxid int `gorm:"column:notagirtaxid;primaryKey" json:"notagirtaxid"`
	Notagirid int `gorm:"column:notagirid" json:"notagirid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Notagirtax) TableName() string {
	return "notagirtax"
}
