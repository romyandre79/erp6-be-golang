package models

import (
  "time"
)

type Snro struct {
	Snroid int `gorm:"column:snroid;primaryKey" json:"snroid"`
	Description string `gorm:"column:description" json:"description"`
	Formatdoc string `gorm:"column:formatdoc" json:"formatdoc"`
	Formatno string `gorm:"column:formatno" json:"formatno"`
	Repeatby *string `gorm:"column:repeatby" json:"repeatby"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Snro) TableName() string {
	return "snro"
}
