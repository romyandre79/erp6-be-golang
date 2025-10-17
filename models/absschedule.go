package models

import (
  "time"
)

type Absschedule struct {
	Absscheduleid int `gorm:"column:absscheduleid;primaryKey" json:"absscheduleid"`
	Absschedulename string `gorm:"column:absschedulename" json:"absschedulename"`
	Absin string `gorm:"column:absin" json:"absin"`
	Absout string `gorm:"column:absout" json:"absout"`
	Absstatusid int `gorm:"column:absstatusid" json:"absstatusid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Absstatus *Absstatus `gorm:"foreignKey:absstatusid;references:absstatusid"`
}

func (Absschedule) TableName() string {
	return "absschedule"
}
