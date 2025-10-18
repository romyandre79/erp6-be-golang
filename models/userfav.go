package models

import (
  "time"
)

type Userfav struct {
	Userfavid int `gorm:"column:userfavid;primaryKey" json:"userfavid"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Menuaccessid int `gorm:"column:menuaccessid" json:"menuaccessid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Menuaccess Menuaccess `gorm:"foreignKey:menuaccessid;references:menuaccessid"`
	Useraccess Useraccess `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Userfav) TableName() string {
	return "userfav"
}
