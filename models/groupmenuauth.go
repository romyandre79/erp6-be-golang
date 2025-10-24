package models

import (
  "time"
)

type Groupmenuauth struct {
	Groupmenuauthid int `gorm:"column:groupmenuauthid;primaryKey" json:"groupmenuauthid"`
	Groupaccessid int `gorm:"column:groupaccessid" json:"groupaccessid"`
	Menuauthid int `gorm:"column:menuauthid" json:"menuauthid"`
	Menuvalueid string `gorm:"column:menuvalueid" json:"menuvalueid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Groupmenuauth) TableName() string {
	return "groupmenuauth"
}
