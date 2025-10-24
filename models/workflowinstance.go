package models

type Workflowinstance struct {
	Workflowinstanceid int    `gorm:"column:workflowinstanceid;primaryKey" json:"workflowinstanceid"`
	Workflowid         int    `gorm:"column:workflowid" json:"workflowid"`
	Datainput          string `gorm:"column:datainput" json:"datainput"`
}

func (Workflowinstance) TableName() string {
	return "workflowinstance"
}
