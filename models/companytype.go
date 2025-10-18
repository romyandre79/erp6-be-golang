package models

import (
  "time"
)

type Companytype struct {
	Companytypeid int `gorm:"column:companytypeid;primaryKey" json:"companytypeid"`
	Companytypename string `gorm:"column:companytypename" json:"companytypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Companytype) TableName() string {
	return "companytype"
}
