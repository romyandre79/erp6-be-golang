package models

import (
  "time"
)

type Contacttype struct {
	Contacttypeid int `gorm:"column:contacttypeid;primaryKey" json:"contacttypeid"`
	Contacttypename string `gorm:"column:contacttypename" json:"contacttypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Contacttype) TableName() string {
	return "contacttype"
}
