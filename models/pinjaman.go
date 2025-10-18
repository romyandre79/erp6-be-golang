package models

import (
  "time"
)

type Pinjaman struct {
	Pinjamanid int `gorm:"column:pinjamanid;primaryKey" json:"pinjamanid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Pinjamandate time.Time `gorm:"column:pinjamandate" json:"pinjamandate"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Barangid int `gorm:"column:barangid" json:"barangid"`
	Pinjamanno *string `gorm:"column:pinjamanno" json:"pinjamanno"`
	Jenispinjamanid int `gorm:"column:jenispinjamanid" json:"jenispinjamanid"`
	Lamaangsuran int `gorm:"column:lamaangsuran" json:"lamaangsuran"`
	Angsuranperbulan *int `gorm:"column:angsuranperbulan" json:"angsuranperbulan"`
	Noperjanjian string `gorm:"column:noperjanjian" json:"noperjanjian"`
	Norekening string `gorm:"column:norekening" json:"norekening"`
	Plafondpinjaman float64 `gorm:"column:plafondpinjaman" json:"plafondpinjaman"`
	Bunga float64 `gorm:"column:bunga" json:"bunga"`
	Biayaadm float64 `gorm:"column:biayaadm" json:"biayaadm"`
	Biayamaterai float64 `gorm:"column:biayamaterai" json:"biayamaterai"`
	Simpananpokok float64 `gorm:"column:simpananpokok" json:"simpananpokok"`
	Simpananwajib float64 `gorm:"column:simpananwajib" json:"simpananwajib"`
	Sisapokok float64 `gorm:"column:sisapokok" json:"sisapokok"`
	Sisabunga float64 `gorm:"column:sisabunga" json:"sisabunga"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Jenispinjaman Jenispinjaman `gorm:"foreignKey:jenispinjamanid;references:jenispinjamanid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Pinjaman) TableName() string {
	return "pinjaman"
}
