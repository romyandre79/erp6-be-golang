package models

import (
  "time"
)

type Onleavetype struct {
	Onleavetypeid int `gorm:"column:onleavetypeid;primaryKey" json:"onleavetypeid"`
	Onleavename string `gorm:"column:onleavename" json:"onleavename"`
	Cutimax int `gorm:"column:cutimax" json:"cutimax"`
	Cutistart int `gorm:"column:cutistart" json:"cutistart"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Absstatusid int `gorm:"column:absstatusid" json:"absstatusid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Onleavetype) TableName() string {
	return "onleavetype"
}
