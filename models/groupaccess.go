package models

import "time"

type Groupaccess struct {
	Groupaccessid int         `gorm:"column:groupaccessid;primaryKey" json:"groupaccessid"`
	Groupname     string      `gorm:"column:groupname" json:"groupname"`
	Description   string      `gorm:"column:description" json:"description"`
	Recordstatus  int         `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Usergroups    []Usergroup `gorm:"foreignKey:Groupaccessid;references:Groupaccessid" json:"usergroups"`
	Groupmenus    []Groupmenu `gorm:"foreignKey:Groupaccessid;references:Groupaccessid" json:"groupmenus"`
}

func (Groupaccess) TableName() string {
	return "groupaccess"
}
