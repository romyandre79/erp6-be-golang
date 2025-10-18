package models

import (
  "time"
)

type Holiday struct {
	Holidayid int `gorm:"column:holidayid;primaryKey" json:"holidayid"`
	Holidaydate time.Time `gorm:"column:holidaydate" json:"holidaydate"`
	Holidayname string `gorm:"column:holidayname" json:"holidayname"`
	Iscuti int8 `gorm:"column:iscuti" json:"iscuti"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Holiday) TableName() string {
	return "holiday"
}
