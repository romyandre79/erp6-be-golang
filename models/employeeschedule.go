package models

import (
  "time"
)

type Employeeschedule struct {
	Employeescheduleid int `gorm:"column:employeescheduleid;primaryKey" json:"employeescheduleid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Scheduledate time.Time `gorm:"column:scheduledate" json:"scheduledate"`
	Absscheduleid int `gorm:"column:absscheduleid" json:"absscheduleid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname *string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Absschedule Absschedule `gorm:"foreignKey:absscheduleid;references:absscheduleid"`
}

func (Employeeschedule) TableName() string {
	return "employeeschedule"
}
