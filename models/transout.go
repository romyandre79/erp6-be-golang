package models

import (
  "time"
)

type Transout struct {
	Transoutid int `gorm:"column:transoutid;primaryKey" json:"transoutid"`
	Docdate time.Time `gorm:"column:docdate" json:"docdate"`
	Transoutno *string `gorm:"column:transoutno" json:"transoutno"`
	Employeeid int `gorm:"column:employeeid" json:"employeeid"`
	Permitexitid int `gorm:"column:permitexitid" json:"permitexitid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Description *string `gorm:"column:description" json:"description"`
	Statusname *string `gorm:"column:statusname" json:"statusname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee Employee `gorm:"foreignKey:employeeid;references:employeeid"`
	Permitexit Permitexit `gorm:"foreignKey:permitexitid;references:permitexitid"`
}

func (Transout) TableName() string {
	return "transout"
}
