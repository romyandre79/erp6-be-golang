package models

import (
  "time"
)

type Booksubject struct {
	Booksubjectid int `gorm:"column:booksubjectid;primaryKey" json:"booksubjectid"`
	Subjectname string `gorm:"column:subjectname" json:"subjectname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Booksubject) TableName() string {
	return "booksubject"
}
