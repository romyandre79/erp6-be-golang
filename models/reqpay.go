package models

import (
  "time"
)

type Reqpay struct {
	Reqpayid int `gorm:"column:reqpayid;primaryKey" json:"reqpayid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Reqpaydate time.Time `gorm:"column:reqpaydate" json:"reqpaydate"`
	Reqpayno *string `gorm:"column:reqpayno" json:"reqpayno"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Iscbout int8 `gorm:"column:iscbout" json:"iscbout"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Reqpay) TableName() string {
	return "reqpay"
}
