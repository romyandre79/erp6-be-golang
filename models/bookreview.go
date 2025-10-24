package models

import (
  "time"
)

type Bookreview struct {
	Bookreviewid int `gorm:"column:bookreviewid;primaryKey" json:"bookreviewid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Reviewnote string `gorm:"column:reviewnote" json:"reviewnote"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Postdate time.Time `gorm:"column:postdate" json:"postdate"`
	Book Book `gorm:"foreignKey:bookid;references:bookid"`
	Useraccess Useraccess `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Bookreview) TableName() string {
	return "bookreview"
}
