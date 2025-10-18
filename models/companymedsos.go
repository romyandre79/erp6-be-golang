package models

import (

)

type Companymedsos struct {
	Companymedsosid int `gorm:"column:companymedsosid;primaryKey" json:"companymedsosid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Mediasocialid int `gorm:"column:mediasocialid" json:"mediasocialid"`
	Mediaaccount string `gorm:"column:mediaaccount" json:"mediaaccount"`
	Mediasocial Mediasocial `gorm:"foreignKey:mediasocialid;references:mediasocialid"`
}

func (Companymedsos) TableName() string {
	return "companymedsos"
}
