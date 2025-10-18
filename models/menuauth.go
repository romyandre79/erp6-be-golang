package models

import (
  "time"
)

type Menuauth struct {
	Menuauthid int `gorm:"column:menuauthid;primaryKey" json:"menuauthid"`
	Menuobject string `gorm:"column:menuobject" json:"menuobject"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Menuauth) TableName() string {
	return "menuauth"
}
