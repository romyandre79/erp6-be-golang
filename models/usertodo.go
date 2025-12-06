package models

import "time"

type Usertodo struct {
	Usertodoid   int       `gorm:"column:usertodoid;primaryKey;autoIncrement" json:"usertodoid"`
	Useraccessid int       `gorm:"column:useraccessid" json:"useraccessid"`
	Tododate     time.Time `gorm:"column:tododate;default:CURRENT_TIMESTAMP" json:"tododate"`
	Menuname     string    `gorm:"column:menuname" json:"menuname"`
	Docno        string    `gorm:"column:docno" json:"docno"`
	Description  string    `gorm:"column:description" json:"description"`
	Isread       int       `gorm:"column:isread" json:"isread"`
	Recordstatus int       `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Usertodo) TableName() string {
	return "usertodo"
}
