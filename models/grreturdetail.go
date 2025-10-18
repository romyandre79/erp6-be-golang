package models

import (
  "time"
)

type Grreturdetail struct {
	Grreturdetailid int `gorm:"column:grreturdetailid;primaryKey" json:"grreturdetailid"`
	Grreturid int `gorm:"column:grreturid" json:"grreturid"`
	Grheaderid int `gorm:"column:grheaderid" json:"grheaderid"`
	Podetailid int `gorm:"column:podetailid" json:"podetailid"`
	Productid *int `gorm:"column:productid" json:"productid"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 float64 `gorm:"column:qty3" json:"qty3"`
	Qty4 float64 `gorm:"column:qty4" json:"qty4"`
	Uomid *int `gorm:"column:uomid" json:"uomid"`
	Uom2id *int `gorm:"column:uom2id" json:"uom2id"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Uom4id *int `gorm:"column:uom4id" json:"uom4id"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Storagebinid int `gorm:"column:storagebinid" json:"storagebinid"`
	Sjsupplier string `gorm:"column:sjsupplier" json:"sjsupplier"`
	Grdetailid int `gorm:"column:grdetailid" json:"grdetailid"`
	Itemnote *string `gorm:"column:itemnote" json:"itemnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Podetail Podetail `gorm:"foreignKey:podetailid;references:podetailid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Storagebin Storagebin `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Sloc Sloc `gorm:"foreignKey:slocid;references:slocid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Grreturdetail) TableName() string {
	return "grreturdetail"
}
