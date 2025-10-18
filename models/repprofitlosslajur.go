package models

import (

)

type Repprofitlosslajur struct {
	Companyid *int `gorm:"column:companyid" json:"companyid"`
	Tahun *string `gorm:"column:tahun" json:"tahun"`
	Bulan *int8 `gorm:"column:bulan" json:"bulan"`
	Nourut *int `gorm:"column:nourut" json:"nourut"`
	Keterangan *string `gorm:"column:keterangan" json:"keterangan"`
	Accountcode *string `gorm:"column:accountcode" json:"accountcode"`
	Plantcode *string `gorm:"column:plantcode" json:"plantcode"`
	Plantvalue *float64 `gorm:"column:plantvalue" json:"plantvalue"`
	Iscount *int8 `gorm:"column:iscount" json:"iscount"`
	Isbold *int8 `gorm:"column:isbold" json:"isbold"`
}

func (Repprofitlosslajur) TableName() string {
	return "repprofitlosslajur"
}
