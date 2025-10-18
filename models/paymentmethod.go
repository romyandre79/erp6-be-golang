package models

import (
  "time"
)

type Paymentmethod struct {
	Paymentmethodid int `gorm:"column:paymentmethodid;primaryKey" json:"paymentmethodid"`
	Paycode string `gorm:"column:paycode" json:"paycode"`
	Paydays int `gorm:"column:paydays" json:"paydays"`
	Paymentname string `gorm:"column:paymentname" json:"paymentname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Paymentmethod) TableName() string {
	return "paymentmethod"
}
