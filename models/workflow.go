package models

import (
	"time"
)

type Workflow struct {
	Workflowid   int       `gorm:"column:workflowid;primaryKey" json:"workflowid"`
	Wfname       string    `gorm:"column:wfname" json:"wfname"`
	Wfdesc       string    `gorm:"column:wfdesc" json:"wfdesc"`
	Wfminstat    int8      `gorm:"column:wfminstat" json:"wfminstat"`
	Wfmaxstat    int8      `gorm:"column:wfmaxstat" json:"wfmaxstat"`
	Flow         string    `gorm:"column:flow" json:"flow"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Workflow) TableName() string {
	return "workflow"
}
