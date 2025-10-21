package models

import (
	"time"
)

type City struct {
	Cityid       int       `gorm:"column:cityid;primaryKey" json:"cityid"`
	Provinceid   int       `gorm:"column:provinceid" json:"provinceid"`
	Citycode     string    `gorm:"column:citycode" json:"citycode"`
	Cityname     string    `gorm:"column:cityname" json:"cityname"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Province     Province  `gorm:"foreignKey:provinceid;references:provinceid"`
}

func (City) TableName() string {
	return "city"
}
