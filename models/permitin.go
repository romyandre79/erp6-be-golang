package models

import (
  "time"
)

type Permitin struct {
	Permitinid int `gorm:"column:permitinid;primaryKey" json:"permitinid"`
	Permitinname string `gorm:"column:permitinname" json:"permitinname"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Snro Snro `gorm:"foreignKey:snroid;references:snroid"`
}

func (Permitin) TableName() string {
	return "permitin"
}
