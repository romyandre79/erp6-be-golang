package models

import (
  "time"
)

type Parameter struct {
	Paramid int `gorm:"column:paramid;primaryKey" json:"paramid"`
	Paramname string `gorm:"column:paramname" json:"paramname"`
	Paramdata string `gorm:"column:paramdata" json:"paramdata"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Parameter) TableName() string {
	return "parameter"
}
