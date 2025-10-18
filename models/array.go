package models

import (

)

type Array struct {
	Nourut int8 `gorm:"column:nourut" json:"nourut"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Dataarray string `gorm:"column:dataarray" json:"dataarray"`
}

func (Array) TableName() string {
	return "array"
}
