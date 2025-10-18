package models

import (
  "time"
)

type Jalurterima struct {
	Jalurterimaid int `gorm:"column:jalurterimaid;primaryKey" json:"jalurterimaid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Jalurname string `gorm:"column:jalurname" json:"jalurname"`
	Persenpagu int `gorm:"column:persenpagu" json:"persenpagu"`
	Tahunajaranid int `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Tahunajaran Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Jalurterima) TableName() string {
	return "jalurterima"
}
