package models

import (

)

type Syaratppdb struct {
	Syaratppdbid int `gorm:"column:syaratppdbid;primaryKey" json:"syaratppdbid"`
	Jalurterimaid int `gorm:"column:jalurterimaid" json:"jalurterimaid"`
	Dokumen string `gorm:"column:dokumen" json:"dokumen"`
	Isrequired int8 `gorm:"column:isrequired" json:"isrequired"`
	Jalurterima Jalurterima `gorm:"foreignKey:jalurterimaid;references:jalurterimaid"`
}

func (Syaratppdb) TableName() string {
	return "syaratppdb"
}
