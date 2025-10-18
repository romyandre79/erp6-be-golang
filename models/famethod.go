package models

import (
  "time"
)

type Famethod struct {
	Famethodid int `gorm:"column:famethodid;primaryKey" json:"famethodid"`
	Methodname string `gorm:"column:methodname" json:"methodname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Famethod) TableName() string {
	return "famethod"
}
