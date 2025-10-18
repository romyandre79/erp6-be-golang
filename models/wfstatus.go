package models

import (
  "time"
)

type Wfstatus struct {
	Wfstatusid int `gorm:"column:wfstatusid;primaryKey" json:"wfstatusid"`
	Workflowid int `gorm:"column:workflowid" json:"workflowid"`
	Wfstat int8 `gorm:"column:wfstat" json:"wfstat"`
	Wfstatusname string `gorm:"column:wfstatusname" json:"wfstatusname"`
	Backcolor string `gorm:"column:backcolor" json:"backcolor"`
	Fontcolor string `gorm:"column:fontcolor" json:"fontcolor"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Wfstatus) TableName() string {
	return "wfstatus"
}
