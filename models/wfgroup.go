package models

import (
  "time"
)

type Wfgroup struct {
	Wfgroupid int `gorm:"column:wfgroupid;primaryKey" json:"wfgroupid"`
	Workflowid int `gorm:"column:workflowid" json:"workflowid"`
	Groupaccessid int `gorm:"column:groupaccessid" json:"groupaccessid"`
	Wfbefstat int8 `gorm:"column:wfbefstat" json:"wfbefstat"`
	Wfrecstat int8 `gorm:"column:wfrecstat" json:"wfrecstat"`
	Sla int `gorm:"column:sla" json:"sla"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Groupaccess Groupaccess `gorm:"foreignKey:groupaccessid;references:groupaccessid"`
	Workflow Workflow `gorm:"foreignKey:workflowid;references:workflowid"`
}

func (Wfgroup) TableName() string {
	return "wfgroup"
}
