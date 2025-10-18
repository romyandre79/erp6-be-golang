package models

import (
  "time"
)

type Salesarea struct {
	Salesareaid int `gorm:"column:salesareaid;primaryKey" json:"salesareaid"`
	Areaname string `gorm:"column:areaname" json:"areaname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Salesarea) TableName() string {
	return "salesarea"
}
