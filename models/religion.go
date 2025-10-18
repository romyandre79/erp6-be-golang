package models

import (
  "time"
)

type Religion struct {
	Religionid int `gorm:"column:religionid;primaryKey" json:"religionid"`
	Religionname string `gorm:"column:religionname" json:"religionname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Religion) TableName() string {
	return "religion"
}
