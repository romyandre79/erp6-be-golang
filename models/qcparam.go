package models

import (
  "time"
)

type Qcparam struct {
	Qcparamid int `gorm:"column:qcparamid;primaryKey" json:"qcparamid"`
	Qcparamname string `gorm:"column:qcparamname" json:"qcparamname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Qcparam) TableName() string {
	return "qcparam"
}
