package models

import (
	"time"
)

type Setoran struct {
	Setoranid       int       `gorm:"column:setoranid;primaryKey" json:"setoranid"`
	Plantid         int       `gorm:"column:plantid" json:"plantid"`
	Setoranno       string    `gorm:"column:setoranno" json:"setoranno"`
	Setorandate     time.Time `gorm:"column:setorandate" json:"setorandate"`
	Simpananid      int       `gorm:"column:simpananid" json:"simpananid"`
	Addressbookid   int       `gorm:"column:addressbookid" json:"addressbookid"`
	Jenissimpananid int       `gorm:"column:jenissimpananid" json:"jenissimpananid"`
	Jumlah          float64   `gorm:"column:jumlah" json:"jumlah"`
	Bunga           float64   `gorm:"column:bunga" json:"bunga"`
	Namasetor       string    `gorm:"column:namasetor" json:"namasetor"`
	Noktpsetor      string    `gorm:"column:noktpsetor" json:"noktpsetor"`
	Notelpsetor     string    `gorm:"column:notelpsetor" json:"notelpsetor"`
	Keterangan      string    `gorm:"column:keterangan" json:"keterangan"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname      string    `gorm:"column:statusname" json:"statusname"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Setoran) TableName() string {
	return "setoran"
}
