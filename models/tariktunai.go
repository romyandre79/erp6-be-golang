package models

import (
	"time"
)

type Tariktunai struct {
	Tariktunaiid    int       `gorm:"column:tariktunaiid;primaryKey" json:"tariktunaiid"`
	Tariktunaidate  time.Time `gorm:"column:tariktunaidate" json:"tariktunaidate"`
	Tariktunaino    string    `gorm:"column:tariktunaino" json:"tariktunaino"`
	Simpananid      int       `gorm:"column:simpananid" json:"simpananid"`
	Plantid         int       `gorm:"column:plantid" json:"plantid"`
	Addressbookid   int       `gorm:"column:addressbookid" json:"addressbookid"`
	Jenissimpananid int       `gorm:"column:jenissimpananid" json:"jenissimpananid"`
	Jumlah          float64   `gorm:"column:jumlah" json:"jumlah"`
	Bunga           float64   `gorm:"column:bunga" json:"bunga"`
	Namatarik       string    `gorm:"column:namatarik" json:"namatarik"`
	Noktptarik      string    `gorm:"column:noktptarik" json:"noktptarik"`
	Notelptarik     string    `gorm:"column:notelptarik" json:"notelptarik"`
	Keterangan      string    `gorm:"column:keterangan" json:"keterangan"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string    `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Tariktunai) TableName() string {
	return "tariktunai"
}
