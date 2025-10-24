package models

import (
  "time"
)

type Jadwalppdb struct {
	Jadwalppdbid int `gorm:"column:jadwalppdbid;primaryKey" json:"jadwalppdbid"`
	Jalurterimaid int `gorm:"column:jalurterimaid" json:"jalurterimaid"`
	Description string `gorm:"column:description" json:"description"`
	Lokasi string `gorm:"column:lokasi" json:"lokasi"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Notes string `gorm:"column:notes" json:"notes"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Jalurterima Jalurterima `gorm:"foreignKey:jalurterimaid;references:jalurterimaid"`
}

func (Jadwalppdb) TableName() string {
	return "jadwalppdb"
}
