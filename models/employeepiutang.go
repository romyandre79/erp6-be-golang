package models

import (
	"time"
)

type Employeepiutang struct {
	Employeepiutangid   int         `gorm:"column:employeepiutangid;primaryKey" json:"employeepiutangid"`
	Employeepiutangno   string      `gorm:"column:employeepiutangno" json:"employeepiutangno"`
	Employeepiutangdate time.Time   `gorm:"column:employeepiutangdate" json:"employeepiutangdate"`
	Duedate             time.Time   `gorm:"column:duedate" json:"duedate"`
	Plantid             int         `gorm:"column:plantid" json:"plantid"`
	Addressbookid       int         `gorm:"column:addressbookid" json:"addressbookid"`
	Positionname        string      `gorm:"column:positionname" json:"positionname"`
	Nilai               float64     `gorm:"column:nilai" json:"nilai"`
	Description         string      `gorm:"column:description" json:"description"`
	Recordstatus        int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname          string      `gorm:"column:statusname" json:"statusname"`
	Updatedate          time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook         Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant               Plant       `gorm:"foreignKey:plantid;references:plantid"`
}

func (Employeepiutang) TableName() string {
	return "employeepiutang"
}
