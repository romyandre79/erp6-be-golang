package models

import (

)

type Repprofitlosstrial struct {
	Companyid *int `gorm:"column:companyid" json:"companyid"`
	Bulan *int `gorm:"column:bulan" json:"bulan"`
	Tahun *int `gorm:"column:tahun" json:"tahun"`
	Accountcode *string `gorm:"column:accountcode" json:"accountcode"`
	Keterangan *string `gorm:"column:keterangan" json:"keterangan"`
	Accountvalue *float64 `gorm:"column:accountvalue" json:"accountvalue"`
	Nourut *int `gorm:"column:nourut" json:"nourut"`
	Plantcode *string `gorm:"column:plantcode" json:"plantcode"`
}

func (Repprofitlosstrial) TableName() string {
	return "repprofitlosstrial"
}
