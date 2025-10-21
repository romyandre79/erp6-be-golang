package models

import (
	"time"
)

type Grdetail struct {
	Grdetailid   int           `gorm:"column:grdetailid;primaryKey" json:"grdetailid"`
	Grheaderid   int           `gorm:"column:grheaderid" json:"grheaderid"`
	Podetailid   int           `gorm:"column:podetailid" json:"podetailid"`
	Productcode  string        `gorm:"column:productcode" json:"productcode"`
	Productname  string        `gorm:"column:productname" json:"productname"`
	Productid    int           `gorm:"column:productid" json:"productid"`
	Qty          float64       `gorm:"column:qty" json:"qty"`
	Qty2         float64       `gorm:"column:qty2" json:"qty2"`
	Qty3         float64       `gorm:"column:qty3" json:"qty3"`
	Uomid        int           `gorm:"column:uomid" json:"uomid"`
	Uom2id       int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id       int           `gorm:"column:uom3id" json:"uom3id"`
	Slocid       int           `gorm:"column:slocid" json:"slocid"`
	Storagebinid int           `gorm:"column:storagebinid" json:"storagebinid"`
	Qtyretur     float64       `gorm:"column:qtyretur" json:"qtyretur"`
	Qty2retur    float64       `gorm:"column:qty2retur" json:"qty2retur"`
	Lotno        string        `gorm:"column:lotno" json:"lotno"`
	Itemnote     string        `gorm:"column:itemnote" json:"itemnote"`
	Updatedate   time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Podetail     Podetail      `gorm:"foreignKey:podetailid;references:podetailid"`
	Product      Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc         Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom          Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2         Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3         Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Grdetail) TableName() string {
	return "grdetail"
}
