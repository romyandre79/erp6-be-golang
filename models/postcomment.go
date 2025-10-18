package models

import (
  "time"
)

type Postcomment struct {
	Postcommentid int `gorm:"column:postcommentid;primaryKey" json:"postcommentid"`
	Postid int `gorm:"column:postid" json:"postid"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Comment string `gorm:"column:comment" json:"comment"`
	Commentdate time.Time `gorm:"column:commentdate" json:"commentdate"`
	Userrate int `gorm:"column:userrate" json:"userrate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Post *Post `gorm:"foreignKey:postcommentid;references:postid"`
	Useraccess *Useraccess `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Postcomment) TableName() string {
	return "postcomment"
}
