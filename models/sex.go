package models

import (
  "time"
)

type Sex struct {
	Sexid int `gorm:"column:sexid;primaryKey" json:"sexid"`
	Sexname string `gorm:"column:sexname" json:"sexname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Sex) TableName() string {
	return "sex"
}
