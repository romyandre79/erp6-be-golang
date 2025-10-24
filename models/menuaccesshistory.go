package models

type Menuaccesshistory struct {
	Menuaccesshistoryid int    `gorm:"column:menuaccesshistoryid;primaryKey" json:"menuaccesshistoryid"`
	Menuaccessid        int    `gorm:"column:menuaccessid" json:"menuaccessid"`
	Menuformdetail      string `gorm:"column:menuformdetail" json:"menuformdetail"`
	Menuform            string `gorm:"column:menuform" json:"menuform"`
}

func (Menuaccesshistory) TableName() string {
	return "menuaccesshistory"
}
