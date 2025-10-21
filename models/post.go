package models

import (
	"time"
)

type Post struct {
	Postid       int       `gorm:"column:postid;primaryKey" json:"postid"`
	Companyid    int       `gorm:"column:companyid" json:"companyid"`
	Useraccessid int       `gorm:"column:useraccessid" json:"useraccessid"`
	Postpic      string    `gorm:"column:postpic" json:"postpic"`
	Title        string    `gorm:"column:title" json:"title"`
	Description  string    `gorm:"column:description" json:"description"`
	Metatag      string    `gorm:"column:metatag" json:"metatag"`
	Slug         string    `gorm:"column:slug" json:"slug"`
	Postupdate   time.Time `gorm:"column:postupdate" json:"postupdate"`
	Created      time.Time `gorm:"column:created" json:"created"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Company      Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Post) TableName() string {
	return "post"
}
