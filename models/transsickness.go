package models

import (
	"time"
)

type Transsickness struct {
	Transsicknessid int       `gorm:"column:transsicknessid;primaryKey" json:"transsicknessid"`
	Docdate         time.Time `gorm:"column:docdate" json:"docdate"`
	Transsickno     string    `gorm:"column:transsickno" json:"transsickno"`
	Employeeid      int       `gorm:"column:employeeid" json:"employeeid"`
	Startdate       time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate         time.Time `gorm:"column:enddate" json:"enddate"`
	Description     string    `gorm:"column:description" json:"description"`
	Suratdokter     string    `gorm:"column:suratdokter" json:"suratdokter"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string    `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee        Employee  `gorm:"foreignKey:employeeid;references:employeeid"`
}

func (Transsickness) TableName() string {
	return "transsickness"
}
