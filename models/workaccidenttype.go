package models

import (
  "time"
)

type Workaccidenttype struct {
	Workaccidenttypeid int `gorm:"column:workaccidenttypeid;primaryKey" json:"workaccidenttypeid"`
	Workaccidenttypename string `gorm:"column:workaccidenttypename" json:"workaccidenttypename"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Workaccidenttype) TableName() string {
	return "workaccidenttype"
}
