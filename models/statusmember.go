package models

import (
  "time"
)

type Statusmember struct {
	Statusmemberid int `gorm:"column:statusmemberid;primaryKey" json:"statusmemberid"`
	Statusmembername string `gorm:"column:statusmembername" json:"statusmembername"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Statusmember) TableName() string {
	return "statusmember"
}
