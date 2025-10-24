package models

import (
  "time"
)

type Processprd struct {
	Processprdid int `gorm:"column:processprdid;primaryKey" json:"processprdid"`
	Processprdname string `gorm:"column:processprdname" json:"processprdname"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Snro Snro `gorm:"foreignKey:snroid;references:snroid"`
}

func (Processprd) TableName() string {
	return "processprd"
}
