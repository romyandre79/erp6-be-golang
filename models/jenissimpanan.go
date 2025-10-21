package models

import (
	"time"
)

type Jenissimpanan struct {
	Jenissimpananid int       `gorm:"column:jenissimpananid;primaryKey" json:"jenissimpananid"`
	Companyid       int       `gorm:"column:companyid" json:"companyid"`
	Namasimpanan    string    `gorm:"column:namasimpanan" json:"namasimpanan"`
	Jumlah          float64   `gorm:"column:jumlah" json:"jumlah"`
	Bunga           float64   `gorm:"column:bunga" json:"bunga"`
	Fixed           int8      `gorm:"column:fixed" json:"fixed"`
	Tenor           int8      `gorm:"column:tenor" json:"tenor"`
	Isauto          int8      `gorm:"column:isauto" json:"isauto"`
	Accsimpanan     int       `gorm:"column:accsimpanan" json:"accsimpanan"`
	Acckas          int       `gorm:"column:acckas" json:"acckas"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company         Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Jenissimpanan) TableName() string {
	return "jenissimpanan"
}
