package models

type Workflowdetail struct {
	Workflowdetailid  int             `gorm:"column:workflowdetailid;primaryKey" json:"workflowdetailid"`
	Componentid       int             `gorm:"column:componentid" json:"componentid"`
	Componentdetailid int             `gorm:"column:componentdetailid" json:"componentdetailid"`
	Workflowid        int             `gorm:"column:workflowid" json:"workflowid"`
	Componentvalue    string          `gorm:"column:componentvalue" json:"componentvalue"`
	Nodeid            int             `gorm:"column:nodeid" json:"nodeid"`
	Componentdetail   Componentdetail `gorm:"foreignKey:componentdetailid;references:componentdetailid"`
	Component         Component       `gorm:"foreignKey:componentid;references:componentid"`
	Workflow          Workflow        `gorm:"foreignKey:workflowid;references:workflowid"`
}

func (Workflowdetail) TableName() string {
	return "workflowdetail"
}
