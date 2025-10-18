package models

import (
  "time"
)

type Addresstype struct {
	Addresstypeid int `gorm:"column:addresstypeid;primaryKey" json:"addresstypeid"`
	Addresstypename string `gorm:"column:addresstypename" json:"addresstypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Addresstype) TableName() string {
	return "addresstype"
}
