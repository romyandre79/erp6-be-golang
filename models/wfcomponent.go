package models

type Wfcomponent struct {
	Wfcomponentid     int    `gorm:"column:wfcomponentid;primaryKey" json:"wfcomponentid"`
	Workflowid        int    `gorm:"column:workflowid" json:"workflowid"`
	Componentname     string `gorm:"column:componentname" json:"componentname"`
	Componentid       int    `gorm:"column:componentid" json:"componentid"`
	Componentdetailid int    `gorm:"column:componentdetailid" json:"componentdetailid"`
	Wfvalue           string `gorm:"column:wfvalue" json:"wfvalue"`
}

func (Wfcomponent) TableName() string {
	return "wfcomponent"
}
