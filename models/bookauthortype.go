package models

import (
  "time"
)

type Bookauthortype struct {
	Bookauthortypeid int `gorm:"column:bookauthortypeid;primaryKey" json:"bookauthortypeid"`
	Authortypename string `gorm:"column:authortypename" json:"authortypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookauthortype) TableName() string {
	return "bookauthortype"
}
