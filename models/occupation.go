package models

import (
  "time"
)

type Occupation struct {
	Occupationid int `gorm:"column:occupationid;primaryKey" json:"occupationid"`
	Occupationname string `gorm:"column:occupationname" json:"occupationname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Occupation) TableName() string {
	return "occupation"
}
