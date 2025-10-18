package models

import (
  "time"
)

type Purchinforec struct {
	Purchinforecid int `gorm:"column:purchinforecid;primaryKey" json:"purchinforecid"`
	Poheaderid *int `gorm:"column:poheaderid" json:"poheaderid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Toleransidown *float64 `gorm:"column:toleransidown" json:"toleransidown"`
	Toleransiup *float64 `gorm:"column:toleransiup" json:"toleransiup"`
	Price *float64 `gorm:"column:price" json:"price"`
	Currencyid *int `gorm:"column:currencyid" json:"currencyid"`
	Biddate *time.Time `gorm:"column:biddate" json:"biddate"`
	Currencyrate float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
}

func (Purchinforec) TableName() string {
	return "purchinforec"
}
