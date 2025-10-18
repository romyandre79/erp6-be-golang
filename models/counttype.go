package models

import (
  "time"
)

type Counttype struct {
	Counttypeid int `gorm:"column:counttypeid;primaryKey" json:"counttypeid"`
	Counttypename string `gorm:"column:counttypename" json:"counttypename"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Counttype) TableName() string {
	return "counttype"
}
