package models

import (
  "time"
)

type Jenisanggota struct {
	Jenisanggotaid int `gorm:"column:jenisanggotaid;primaryKey" json:"jenisanggotaid"`
	Kodejenisanggota string `gorm:"column:kodejenisanggota" json:"kodejenisanggota"`
	Namajenisanggota string `gorm:"column:namajenisanggota" json:"namajenisanggota"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Snro Snro `gorm:"foreignKey:snroid;references:snroid"`
}

func (Jenisanggota) TableName() string {
	return "jenisanggota"
}
