package models

import (
  "time"
)

type Userdash struct {
	Userdashid int `gorm:"column:userdashid;primaryKey" json:"userdashid"`
	Groupaccessid int `gorm:"column:groupaccessid" json:"groupaccessid"`
	Widgetid int `gorm:"column:widgetid" json:"widgetid"`
	Menuaccessid int `gorm:"column:menuaccessid" json:"menuaccessid"`
	Position int8 `gorm:"column:position" json:"position"`
	Webformat string `gorm:"column:webformat" json:"webformat"`
	Dashgroup int8 `gorm:"column:dashgroup" json:"dashgroup"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Menuaccess Menuaccess `gorm:"foreignKey:menuaccessid;references:menuaccessid"`
	Widget Widget `gorm:"foreignKey:widgetid;references:widgetid"`
}

func (Userdash) TableName() string {
	return "userdash"
}
