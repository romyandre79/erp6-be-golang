package models

import (
  "time"
)

type Appsmenu struct {
	Appsmenuid int `gorm:"column:appsmenuid;primaryKey" json:"appsmenuid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Appsid int `gorm:"column:appsid" json:"appsid"`
	Menuaccessid int `gorm:"column:menuaccessid" json:"menuaccessid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Apps Apps `gorm:"foreignKey:appsid;references:appsid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Menuaccess Menuaccess `gorm:"foreignKey:menuaccessid;references:menuaccessid"`
}

func (Appsmenu) TableName() string {
	return "appsmenu"
}
