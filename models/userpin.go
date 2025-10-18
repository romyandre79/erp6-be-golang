package models

import (
  "time"
)

type Userpin struct {
	Userpinid int `gorm:"column:userpinid;primaryKey" json:"userpinid"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Menuaccessid int `gorm:"column:menuaccessid" json:"menuaccessid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Userpin) TableName() string {
	return "userpin"
}
