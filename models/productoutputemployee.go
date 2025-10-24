package models

import (
	"time"
)

type Productoutputemployee struct {
	Productoutputemployeeid int       `gorm:"column:productoutputemployeeid;primaryKey" json:"productoutputemployeeid"`
	Productoutputid         int       `gorm:"column:productoutputid" json:"productoutputid"`
	Productoutputfgid       int       `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Employeeid              int       `gorm:"column:employeeid" json:"employeeid"`
	Fullname                string    `gorm:"column:fullname" json:"fullname"`
	Oldnik                  string    `gorm:"column:oldnik" json:"oldnik"`
	Description             string    `gorm:"column:description" json:"description"`
	Updatedate              time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee                Employee  `gorm:"foreignKey:employeeid;references:employeeid"`
}

func (Productoutputemployee) TableName() string {
	return "productoutputemployee"
}
