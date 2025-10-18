package models

import (
  "time"
)

type Recurringtype struct {
	Recurringtypeid int `gorm:"column:recurringtypeid;primaryKey" json:"recurringtypeid"`
	Recurringtypename string `gorm:"column:recurringtypename" json:"recurringtypename"`
	Recurringtypevalue int `gorm:"column:recurringtypevalue" json:"recurringtypevalue"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Recurringtype) TableName() string {
	return "recurringtype"
}
