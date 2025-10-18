package models

import (
  "time"
)

type Fixasset struct {
	Fixassetid int `gorm:"column:fixassetid;primaryKey" json:"fixassetid"`
	Assetno *string `gorm:"column:assetno" json:"assetno"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Pono *string `gorm:"column:pono" json:"pono"`
	Buydate time.Time `gorm:"column:buydate" json:"buydate"`
	Productid int `gorm:"column:productid" json:"productid"`
	Description string `gorm:"column:description" json:"description"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Qty2 int `gorm:"column:qty2" json:"qty2"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Price float64 `gorm:"column:price" json:"price"`
	Nilairesidu float64 `gorm:"column:nilairesidu" json:"nilairesidu"`
	Famethodid int `gorm:"column:famethodid" json:"famethodid"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Umur float64 `gorm:"column:umur" json:"umur"`
	Ratesusut float64 `gorm:"column:ratesusut" json:"ratesusut"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Fixasset) TableName() string {
	return "fixasset"
}
