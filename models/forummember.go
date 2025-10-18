package models

import (

)

type Forummember struct {
	Forummemberid int `gorm:"column:forummemberid;primaryKey" json:"forummemberid"`
	Forumid int `gorm:"column:forumid" json:"forumid"`
	Forum Forum `gorm:"foreignKey:forumid;references:forumid"`
}

func (Forummember) TableName() string {
	return "forummember"
}
