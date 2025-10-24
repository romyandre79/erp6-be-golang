package models

import (
	"time"
)

type Reportperday struct {
	Reportperdayid int         `gorm:"column:reportperdayid;primaryKey" json:"reportperdayid"`
	Addressbookid  int         `gorm:"column:addressbookid" json:"addressbookid"`
	Absdate        time.Time   `gorm:"column:absdate" json:"absdate"`
	Hourin         string      `gorm:"column:hourin" json:"hourin"`
	Hourout        string      `gorm:"column:hourout" json:"hourout"`
	Absscheduleid  int         `gorm:"column:absscheduleid" json:"absscheduleid"`
	Statusin       string      `gorm:"column:statusin" json:"statusin"`
	Statusout      string      `gorm:"column:statusout" json:"statusout"`
	Reason         string      `gorm:"column:reason" json:"reason"`
	Addressbook    Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Reportperday) TableName() string {
	return "reportperday"
}
