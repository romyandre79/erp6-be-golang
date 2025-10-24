package models

import (
	"time"
)

type Postcategory struct {
	Postcategoryid int       `gorm:"column:postcategoryid;primaryKey" json:"postcategoryid"`
	Postid         int       `gorm:"column:postid" json:"postid"`
	Categoryid     int       `gorm:"column:categoryid" json:"categoryid"`
	Updatedate     time.Time `gorm:"column:updatedate" json:"updatedate"`
	Category       Category  `gorm:"foreignKey:categoryid;references:categoryid"`
	Post           Post      `gorm:"foreignKey:postid;references:postid"`
}

func (Postcategory) TableName() string {
	return "postcategory"
}
