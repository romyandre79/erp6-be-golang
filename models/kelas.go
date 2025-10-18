package models

import (
	"time"
)

type Kelas struct {
	Kelasid      int            `gorm:"column:kelasid;primaryKey" json:"kelasid"`
	Companyid    int            `gorm:"column:companyid" json:"companyid"`
	Jurusanid    int            `gorm:"column:jurusanid" json:"jurusanid"`
	Kelasname    string         `gorm:"column:kelasname" json:"kelasname"`
	Urutan       *int8          `gorm:"column:urutan" json:"urutan"`
	Pagu         int8           `gorm:"column:pagu" json:"pagu"`
	Waktustart   *time.Time     `gorm:"column:waktustart" json:"waktustart"`
	Waktuend     *time.Time     `gorm:"column:waktuend" json:"waktuend"`
	Description  *string        `gorm:"column:description" json:"description"`
	Photo        *string        `gorm:"column:photo" json:"photo"`
	Bobot        *int8          `gorm:"column:bobot" json:"bobot"`
	Usiarata     *string        `gorm:"column:usiarata" json:"usiarata"`
	Recordstatus int8           `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate  time.Time      `gorm:"column:createddate" json:"createddate"`
	Updatedate   time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Company      Company        `gorm:"foreignKey:companyid;references:companyid"`
	Jurusan      Educationmajor `gorm:"foreignKey:jurusanid;references:jurusanid"`
}

func (Kelas) TableName() string {
	return "kelas"
}
