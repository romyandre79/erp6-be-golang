package models

import (
  "time"
)

type Groupcustomer struct {
	Groupcustomerid int `gorm:"column:groupcustomerid;primaryKey" json:"groupcustomerid"`
	Groupname string `gorm:"column:groupname" json:"groupname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Groupcustomer) TableName() string {
	return "groupcustomer"
}
