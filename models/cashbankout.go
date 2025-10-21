package models

import (
	"time"
)

type Cashbankout struct {
	Cashbankoutid   int       `gorm:"column:cashbankoutid;primaryKey" json:"cashbankoutid"`
	Plantid         int       `gorm:"column:plantid" json:"plantid"`
	Cashbankoutno   string    `gorm:"column:cashbankoutno" json:"cashbankoutno"`
	Cashbankoutdate time.Time `gorm:"column:cashbankoutdate" json:"cashbankoutdate"`
	Accountid       int       `gorm:"column:accountid" json:"accountid"`
	Reqpayid        int       `gorm:"column:reqpayid" json:"reqpayid"`
	Headernote      string    `gorm:"column:headernote" json:"headernote"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string    `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant           Plant     `gorm:"foreignKey:plantid;references:plantid"`
	Reqpay          Reqpay    `gorm:"foreignKey:reqpayid;references:reqpayid"`
}

func (Cashbankout) TableName() string {
	return "cashbankout"
}
