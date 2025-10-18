package models

import (
  "time"
)

type Levelorg struct {
	Levelorgid int `gorm:"column:levelorgid;primaryKey" json:"levelorgid"`
	Levelorgname string `gorm:"column:levelorgname" json:"levelorgname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Levelorg) TableName() string {
	return "levelorg"
}
