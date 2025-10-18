package models

import (
  "time"
)

type Jenispinjaman struct {
	Jenispinjamanid int `gorm:"column:jenispinjamanid;primaryKey" json:"jenispinjamanid"`
	Namapinjaman string `gorm:"column:namapinjaman" json:"namapinjaman"`
	Jumlah float64 `gorm:"column:jumlah" json:"jumlah"`
	Bunga float64 `gorm:"column:bunga" json:"bunga"`
	Fixed int8 `gorm:"column:fixed" json:"fixed"`
	Biayaadm float64 `gorm:"column:biayaadm" json:"biayaadm"`
	Simpokok float64 `gorm:"column:simpokok" json:"simpokok"`
	Biayamaterai float64 `gorm:"column:biayamaterai" json:"biayamaterai"`
	Biayaasuransi float64 `gorm:"column:biayaasuransi" json:"biayaasuransi"`
	Maxday int8 `gorm:"column:maxday" json:"maxday"`
	Tenor int8 `gorm:"column:tenor" json:"tenor"`
	Isauto int8 `gorm:"column:isauto" json:"isauto"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Jenispinjaman) TableName() string {
	return "jenispinjaman"
}
