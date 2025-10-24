package models

import (
  "time"
)

type Materialstatus struct {
	Materialstatusid int `gorm:"column:materialstatusid;primaryKey" json:"materialstatusid"`
	Materialstatusname string `gorm:"column:materialstatusname" json:"materialstatusname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Materialstatus) TableName() string {
	return "materialstatus"
}
