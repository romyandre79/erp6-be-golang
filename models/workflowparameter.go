package models

import (

)

type Workflowparameter struct {
	Wfparameterid int `gorm:"column:wfparameterid;primaryKey" json:"wfparameterid"`
	Workflowid int `gorm:"column:workflowid" json:"workflowid"`
	Parametername string `gorm:"column:parametername" json:"parametername"`
}

func (Workflowparameter) TableName() string {
	return "workflowparameter"
}
