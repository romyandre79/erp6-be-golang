package models

import (
"time"
)

type Apps struct {
	Appsid int `gorm:"column:appsid;primaryKey" json:"appsid"`
	Appsname string `gorm:"column:appsname" json:"appsname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Apps) TableName() string {
	return "apps"
}
