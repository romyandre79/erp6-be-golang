package models

import "time"

type Usergroup struct {
	Usergroupid   int         `gorm:"column:usergroupid;primaryKey" json:"usergroupid"`
	Useraccessid  int         `gorm:"column:useraccessid" json:"useraccessid"`
	Groupaccessid int         `gorm:"column:groupaccessid" json:"groupaccessid"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Groupaccess   Groupaccess `gorm:"foreignKey:Groupaccessid;references:Groupaccessid"`
	UserAccess    Useraccess  `gorm:"foreignKey:Useraccessid;references:Useraccessid"`
}

func (Usergroup) TableName() string {
	return "usergroup"
}
