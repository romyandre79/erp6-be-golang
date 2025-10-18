package models

import (
  "time"
)

type Kebangsaan struct {
	Kebangsaanid int `gorm:"column:kebangsaanid;primaryKey" json:"kebangsaanid"`
	Kebangsaanname string `gorm:"column:kebangsaanname" json:"kebangsaanname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Kebangsaan) TableName() string {
	return "kebangsaan"
}
