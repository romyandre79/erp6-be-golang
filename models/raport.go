package models

import (

)

type Raport struct {
	Raportid int `gorm:"column:raportid;primaryKey" json:"raportid"`
	Kelasid *int `gorm:"column:kelasid" json:"kelasid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Tahunajaranid *int `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Mapelid *int `gorm:"column:mapelid" json:"mapelid"`
	Nilai *int `gorm:"column:nilai" json:"nilai"`
	Kkm *int `gorm:"column:kkm" json:"kkm"`
}

func (Raport) TableName() string {
	return "raport"
}
