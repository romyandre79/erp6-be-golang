package models

type Workflowparameter struct {
	Wfparameterid int      `gorm:"column:wfparameterid;primaryKey" json:"wfparameterid"`
	Workflowid    int      `gorm:"column:workflowid" json:"workflowid"`
	Parametername string   `gorm:"column:parametername" json:"parametername"`
	Workflow      Workflow `gorm:"foreignKey:Workflowid;references:Workflowid"`
}

func (Workflowparameter) TableName() string {
	return "workflowparameter"
}
