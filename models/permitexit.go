package models

import (
  "time"
)

type Permitexit struct {
	Permitexitid int `gorm:"column:permitexitid;primaryKey" json:"permitexitid"`
	Permitexitname string `gorm:"column:permitexitname" json:"permitexitname"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Snro Snro `gorm:"foreignKey:snroid;references:snroid"`
}

func (Permitexit) TableName() string {
	return "permitexit"
}
