package models

import (
  "time"
)

type Taxpo struct {
	Taxpoid int `gorm:"column:taxpoid;primaryKey" json:"taxpoid"`
	Poheaderid int `gorm:"column:poheaderid" json:"poheaderid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Taxpo) TableName() string {
	return "taxpo"
}
