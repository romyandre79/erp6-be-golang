package models

import (
	"time"
)

type Onleavetrans struct {
	Onleavetransid int         `gorm:"column:onleavetransid;primaryKey" json:"onleavetransid"`
	Docdate        time.Time   `gorm:"column:docdate" json:"docdate"`
	Onleavetransno string      `gorm:"column:onleavetransno" json:"onleavetransno"`
	Employeeid     int         `gorm:"column:employeeid" json:"employeeid"`
	Onleavetypeid  int         `gorm:"column:onleavetypeid" json:"onleavetypeid"`
	Startdate      time.Time   `gorm:"column:startdate" json:"startdate"`
	Enddate        time.Time   `gorm:"column:enddate" json:"enddate"`
	Description    string      `gorm:"column:description" json:"description"`
	Recordstatus   int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string      `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Employee       Employee    `gorm:"foreignKey:employeeid;references:employeeid"`
	Onleavetype    Onleavetype `gorm:"foreignKey:onleavetypeid;references:onleavetypeid"`
}

func (Onleavetrans) TableName() string {
	return "onleavetrans"
}
