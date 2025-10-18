package models

import (
  "time"
)

type Productconversiondetail struct {
	Productconversiondetailid int `gorm:"column:productconversiondetailid;primaryKey" json:"productconversiondetailid"`
	Productconversionid *int `gorm:"column:productconversionid" json:"productconversionid"`
	Productid *int `gorm:"column:productid" json:"productid"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Productconversiondetail) TableName() string {
	return "productconversiondetail"
}
