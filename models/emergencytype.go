package models

import (
  "time"
)

type Emergencytype struct {
	Emergencytypeid int `gorm:"column:emergencytypeid;primaryKey" json:"emergencytypeid"`
	Reporttype string `gorm:"column:reporttype" json:"reporttype"`
	Emergencyname string `gorm:"column:emergencyname" json:"emergencyname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Emergencytype) TableName() string {
	return "emergencytype"
}
