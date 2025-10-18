package models

import (

)

type Companyajaran struct {
	Companyajaranid int `gorm:"column:companyajaranid;primaryKey" json:"companyajaranid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Tahunajaranid int `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Tahunajaran Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Companyajaran) TableName() string {
	return "companyajaran"
}
