package models

import (
  "time"
)

type Filehistory struct {
	Filehistoryid int `gorm:"column:filehistoryid" json:"filehistoryid"`
	Filename string `gorm:"column:filename" json:"filename"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Filecontent string `gorm:"column:filecontent" json:"filecontent"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Filehistory) TableName() string {
	return "filehistory"
}
