package models

import (
  "time"
)

type Chat struct {
	Chatid int `gorm:"column:chatid;primaryKey" json:"chatid"`
	Msgfrom string `gorm:"column:msgfrom" json:"msgfrom"`
	Msgparam *string `gorm:"column:msgparam" json:"msgparam"`
	Msgresponid *int `gorm:"column:msgresponid" json:"msgresponid"`
	Seqorder int8 `gorm:"column:seqorder" json:"seqorder"`
	Parentchatid *int `gorm:"column:parentchatid" json:"parentchatid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Chat) TableName() string {
	return "chat"
}
