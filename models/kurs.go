package models

import (
  "time"
)

type Kurs struct {
	Kursid int `gorm:"column:kursid;primaryKey" json:"kursid"`
	Kursdate time.Time `gorm:"column:kursdate" json:"kursdate"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Taxrate float64 `gorm:"column:taxrate" json:"taxrate"`
	Birate float64 `gorm:"column:birate" json:"birate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Kurs) TableName() string {
	return "kurs"
}
