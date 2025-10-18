package models

import (
  "time"
)

type Requestedby struct {
	Requestedbyid int `gorm:"column:requestedbyid;primaryKey" json:"requestedbyid"`
	Requestedbycode string `gorm:"column:requestedbycode" json:"requestedbycode"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Requestedby) TableName() string {
	return "requestedby"
}
