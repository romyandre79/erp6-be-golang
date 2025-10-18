package models

import (
  "time"
)

type Fahistory struct {
	Fahistoryid int64 `gorm:"column:fahistoryid;primaryKey" json:"fahistoryid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Fixassetid int `gorm:"column:fixassetid" json:"fixassetid"`
	Bulanke int `gorm:"column:bulanke" json:"bulanke"`
	Susutdate time.Time `gorm:"column:susutdate" json:"susutdate"`
	Nilai float64 `gorm:"column:nilai" json:"nilai"`
	Beban float64 `gorm:"column:beban" json:"beban"`
	Nilaiakum float64 `gorm:"column:nilaiakum" json:"nilaiakum"`
	Nilaibuku float64 `gorm:"column:nilaibuku" json:"nilaibuku"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Fahistory) TableName() string {
	return "fahistory"
}
