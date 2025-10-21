package models

import (
	"time"
)

type Pesertappdb struct {
	Pesertappdbid   int       `gorm:"column:pesertappdbid;primaryKey" json:"pesertappdbid"`
	Companyid       int       `gorm:"column:companyid" json:"companyid"`
	Addressbookname string    `gorm:"column:addressbookname" json:"addressbookname"`
	Cityid          int       `gorm:"column:cityid" json:"cityid"`
	Birthdate       time.Time `gorm:"column:birthdate" json:"birthdate"`
	Nisn            string    `gorm:"column:nisn" json:"nisn"`
	Nis             string    `gorm:"column:nis" json:"nis"`
	Useraccessid    int       `gorm:"column:useraccessid" json:"useraccessid"`
	Address         string    `gorm:"column:address" json:"address"`
	Jalurterimaid   int       `gorm:"column:jalurterimaid" json:"jalurterimaid"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate     time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company         Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Pesertappdb) TableName() string {
	return "pesertappdb"
}
