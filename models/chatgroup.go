package models

import (
  "time"
)

type Chatgroup struct {
	Chatgroupid int `gorm:"column:chatgroupid;primaryKey" json:"chatgroupid"`
	Chatid int `gorm:"column:chatid" json:"chatid"`
	Groupaccessid int `gorm:"column:groupaccessid" json:"groupaccessid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Chatgroup) TableName() string {
	return "chatgroup"
}
