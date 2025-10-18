package models

import (
  "time"
)

type Employeejamsostek struct {
	Employeejamsostekid int `gorm:"column:employeejamsostekid;primaryKey" json:"employeejamsostekid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Jamsostekdate *time.Time `gorm:"column:jamsostekdate" json:"jamsostekdate"`
	Jamsostekno string `gorm:"column:jamsostekno" json:"jamsostekno"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeejamsostek) TableName() string {
	return "employeejamsostek"
}
