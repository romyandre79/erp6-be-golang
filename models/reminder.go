package models

import (
  "time"
)

type Reminder struct {
	Reminderid int `gorm:"column:reminderid;primaryKey" json:"reminderid"`
	Activityid int `gorm:"column:activityid" json:"activityid"`
	Remindername string `gorm:"column:remindername" json:"remindername"`
	Description string `gorm:"column:description" json:"description"`
	Reminderdate time.Time `gorm:"column:reminderdate" json:"reminderdate"`
	Repeatformula string `gorm:"column:repeatformula" json:"repeatformula"`
	Repeatmax int `gorm:"column:repeatmax" json:"repeatmax"`
	Repeatcount int `gorm:"column:repeatcount" json:"repeatcount"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Reminder) TableName() string {
	return "reminder"
}
