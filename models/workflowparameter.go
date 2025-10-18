package models

import (

)

type Workflowparameter struct {
	Wfparameterid int `gorm:"column:wfparameterid;primaryKey" json:"wfparameterid"`
	Workflowid int `gorm:"column:workflowid" json:"workflowid"`
	Parametername string `gorm:"column:parametername" json:"parametername"`
	Workflow Workflow `gorm:"foreignKey:workflowid;references:workflowid"`
}

func (Workflowparameter) TableName() string {
	return "workflowparameter"
}
