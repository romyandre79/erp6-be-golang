package models

import (
  "time"
)

type Productsales struct {
	Productsalesid int `gorm:"column:productsalesid;primaryKey" json:"productsalesid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Currencyvalue float64 `gorm:"column:currencyvalue" json:"currencyvalue"`
	Pricecategoryid int `gorm:"column:pricecategoryid" json:"pricecategoryid"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Pricecategory Pricecategory `gorm:"foreignKey:pricecategoryid;references:pricecategoryid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Productsales) TableName() string {
	return "productsales"
}
