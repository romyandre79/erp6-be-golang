package models

import (

)

type Workflowhistory struct {
	Workflowhistoryid int `gorm:"column:workflowhistoryid;primaryKey" json:"workflowhistoryid"`
	Workflowid int `gorm:"column:workflowid" json:"workflowid"`
	Menuflow string `gorm:"column:menuflow" json:"menuflow"`
	Ispublish int8 `gorm:"column:ispublish" json:"ispublish"`
}

func (Workflowhistory) TableName() string {
	return "workflowhistory"
}
