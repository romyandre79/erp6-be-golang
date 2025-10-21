package models

import (
	"time"
)

type Dbobject struct {
	Dbobjectid    int       `gorm:"column:dbobjectid;primaryKey" json:"dbobjectid"`
	Objectname    string    `gorm:"column:objectname" json:"objectname"`
	Objecttype    string    `gorm:"column:objecttype" json:"objecttype"`
	Objectcontent string    `gorm:"column:objectcontent" json:"objectcontent"`
	Objectversion int       `gorm:"column:objectversion" json:"objectversion"`
	Ispublished   int       `gorm:"column:ispublished" json:"ispublished"`
	Sort          int8      `gorm:"column:sort" json:"sort"`
	Comment       string    `gorm:"column:comment" json:"comment"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Dbobject) TableName() string {
	return "dbobject"
}
