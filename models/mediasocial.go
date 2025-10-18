package models

import (

)

type Mediasocial struct {
	Mediasocialid int `gorm:"column:mediasocialid;primaryKey" json:"mediasocialid"`
	Mediasocialname string `gorm:"column:mediasocialname" json:"mediasocialname"`
}

func (Mediasocial) TableName() string {
	return "mediasocial"
}
