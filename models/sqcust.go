package models

import (
	"time"
)

type Sqcust struct {
	Sqcustid      int         `gorm:"column:sqcustid;primaryKey" json:"sqcustid"`
	Sqheaderid    int         `gorm:"column:sqheaderid" json:"sqheaderid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Sqcust) TableName() string {
	return "sqcust"
}
