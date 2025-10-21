package models

import (
	"time"
)

type Mapelkelas struct {
	Mapelkelasid   int         `gorm:"column:mapelkelasid;primaryKey" json:"mapelkelasid"`
	Tahunajaranid  int         `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Mapelid        int         `gorm:"column:mapelid" json:"mapelid"`
	Kelasid        int         `gorm:"column:kelasid" json:"kelasid"`
	Addressbookid  int         `gorm:"column:addressbookid" json:"addressbookid"`
	Mapelvid       string      `gorm:"column:mapelvid" json:"mapelvid"`
	Mapelkelasdesc string      `gorm:"column:mapelkelasdesc" json:"mapelkelasdesc"`
	Createddate    time.Time   `gorm:"column:createddate" json:"createddate"`
	Updatedate     time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook    Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Kelas          Kelas       `gorm:"foreignKey:kelasid;references:kelasid"`
	Mapel          Mapel       `gorm:"foreignKey:mapelid;references:mapelid"`
	Tahunajaran    Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Mapelkelas) TableName() string {
	return "mapelkelas"
}
