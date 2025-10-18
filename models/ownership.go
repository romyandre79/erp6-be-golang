package models

import (
  "time"
)

type Ownership struct {
	Ownershipid int `gorm:"column:ownershipid;primaryKey" json:"ownershipid"`
	Ownershipname string `gorm:"column:ownershipname" json:"ownershipname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Ownership) TableName() string {
	return "ownership"
}
