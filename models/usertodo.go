package models

import (
	"time"
)

type Usertodo struct {
	Usertodoid   int        `gorm:"column:usertodoid;primaryKey" json:"usertodoid"`
	Useraccessid int        `gorm:"column:useraccessid" json:"useraccessid"`
	Tododate     time.Time  `gorm:"column:tododate" json:"tododate"`
	Menuname     string     `gorm:"column:menuname" json:"menuname"`
	Docno        string     `gorm:"column:docno" json:"docno"`
	Description  string     `gorm:"column:description" json:"description"`
	Isread       int8       `gorm:"column:isread" json:"isread"`
	Recordstatus int8       `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time  `gorm:"column:updatedate" json:"updatedate"`
	Useraccess   Useraccess `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Usertodo) TableName() string {
	return "usertodo"
}
