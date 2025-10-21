package models

import (
	"time"
)

type Taxsq struct {
	Taxsqid     int       `gorm:"column:taxsqid;primaryKey" json:"taxsqid"`
	Sqheaderid  int       `gorm:"column:sqheaderid" json:"sqheaderid"`
	Taxid       int       `gorm:"column:taxid" json:"taxid"`
	Description string    `gorm:"column:description" json:"description"`
	Updatedate  time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Taxsq) TableName() string {
	return "taxsq"
}
