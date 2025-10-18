package models

import (
  "time"
)

type Dbobjecthistory struct {
	Dbobjecthistoryid int `gorm:"column:dbobjecthistoryid;primaryKey" json:"dbobjecthistoryid"`
	Dbobjectid int `gorm:"column:dbobjectid" json:"dbobjectid"`
	Objectname string `gorm:"column:objectname" json:"objectname"`
	Objecttype string `gorm:"column:objecttype" json:"objecttype"`
	Objectcontent string `gorm:"column:objectcontent" json:"objectcontent"`
	Objectversion int `gorm:"column:objectversion" json:"objectversion"`
	Createdby string `gorm:"column:createdby" json:"createdby"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Dbobjecthistory) TableName() string {
	return "dbobjecthistory"
}
