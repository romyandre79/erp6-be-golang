package models

import (
	"time"
)

type Grresult struct {
	Grresultid  int           `gorm:"column:grresultid;primaryKey" json:"grresultid"`
	Grheaderid  int           `gorm:"column:grheaderid" json:"grheaderid"`
	Poresultid  int           `gorm:"column:poresultid" json:"poresultid"`
	Productid   int           `gorm:"column:productid" json:"productid"`
	Productcode string        `gorm:"column:productcode" json:"productcode"`
	Productname string        `gorm:"column:productname" json:"productname"`
	Qty1        float64       `gorm:"column:qty1" json:"qty1"`
	Uomid       int           `gorm:"column:uomid" json:"uomid"`
	Qty2        float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id      int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3        float64       `gorm:"column:qty3" json:"qty3"`
	Uom3id      int           `gorm:"column:uom3id" json:"uom3id"`
	Prqty       float64       `gorm:"column:prqty" json:"prqty"`
	Prqty2      float64       `gorm:"column:prqty2" json:"prqty2"`
	Prqty3      float64       `gorm:"column:prqty3" json:"prqty3"`
	Tsqty       float64       `gorm:"column:tsqty" json:"tsqty"`
	Tsqty2      float64       `gorm:"column:tsqty2" json:"tsqty2"`
	Tsqty3      float64       `gorm:"column:tsqty3" json:"tsqty3"`
	Description string        `gorm:"column:description" json:"description"`
	Updatedate  time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product     Product       `gorm:"foreignKey:productid;references:productid"`
	Uom         Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2        Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3        Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Grresult) TableName() string {
	return "grresult"
}
