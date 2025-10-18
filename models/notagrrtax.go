package models

import (
  "time"
)

type Notagrrtax struct {
	Notagrrtaxid int `gorm:"column:notagrrtaxid;primaryKey" json:"notagrrtaxid"`
	Notagrrid int `gorm:"column:notagrrid" json:"notagrrid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Notagrrtax) TableName() string {
	return "notagrrtax"
}
