package models

import (
  "time"
)

type Simpanan struct {
	Simpananid int `gorm:"column:simpananid;primaryKey" json:"simpananid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Simpananno *string `gorm:"column:simpananno" json:"simpananno"`
	Simpanandate time.Time `gorm:"column:simpanandate" json:"simpanandate"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Jenissimpananid int `gorm:"column:jenissimpananid" json:"jenissimpananid"`
	Jumlah float64 `gorm:"column:jumlah" json:"jumlah"`
	Tenor int `gorm:"column:tenor" json:"tenor"`
	Bunga int `gorm:"column:bunga" json:"bunga"`
	Acckas int `gorm:"column:acckas" json:"acckas"`
	Accsimpanan int `gorm:"column:accsimpanan" json:"accsimpanan"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Simpanan) TableName() string {
	return "simpanan"
}
