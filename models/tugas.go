package models

import (
  "time"
)

type Tugas struct {
	Tugasid int `gorm:"column:tugasid;primaryKey" json:"tugasid"`
	Mapelid int `gorm:"column:mapelid" json:"mapelid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Typeid int8 `gorm:"column:typeid" json:"typeid"`
	Judul string `gorm:"column:judul" json:"judul"`
	Durasi int `gorm:"column:durasi" json:"durasi"`
	Info string `gorm:"column:info" json:"info"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Istampilsiswa int8 `gorm:"column:istampilsiswa" json:"istampilsiswa"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
}

func (Tugas) TableName() string {
	return "tugas"
}
