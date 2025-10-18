package models

import (
  "time"
)

type Bookcontenttype struct {
	Bookcontenttypeid int `gorm:"column:bookcontenttypeid;primaryKey" json:"bookcontenttypeid"`
	Contenttypename string `gorm:"column:contenttypename" json:"contenttypename"`
	Contenttypecode string `gorm:"column:contenttypecode" json:"contenttypecode"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookcontenttype) TableName() string {
	return "bookcontenttype"
}
