package models

import (
  "time"
)

type Menuversion struct {
	Menuversionid int `gorm:"column:menuversionid;primaryKey" json:"menuversionid"`
	Menuaccessid int `gorm:"column:menuaccessid" json:"menuaccessid"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Viewcode string `gorm:"column:viewcode" json:"viewcode"`
	Controllercode string `gorm:"column:controllercode" json:"controllercode"`
	Version int `gorm:"column:version" json:"version"`
	Ispublished int8 `gorm:"column:ispublished" json:"ispublished"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Menuaccess Menuaccess `gorm:"foreignKey:menuaccessid;references:menuaccessid"`
	Useraccess Useraccess `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Menuversion) TableName() string {
	return "menuversion"
}
