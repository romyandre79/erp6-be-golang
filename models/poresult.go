package models

import (
	"time"
)

type Poresult struct {
	Poresultid  int           `gorm:"column:poresultid;primaryKey" json:"poresultid"`
	Poheaderid  int           `gorm:"column:poheaderid" json:"poheaderid"`
	Prresultid  int           `gorm:"column:prresultid" json:"prresultid"`
	Prheaderid  int           `gorm:"column:prheaderid" json:"prheaderid"`
	Productid   int           `gorm:"column:productid" json:"productid"`
	Productcode string        `gorm:"column:productcode" json:"productcode"`
	Productname string        `gorm:"column:productname" json:"productname"`
	Qty         float64       `gorm:"column:qty" json:"qty"`
	Uomid       int           `gorm:"column:uomid" json:"uomid"`
	Qty2        float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id      int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3        float64       `gorm:"column:qty3" json:"qty3"`
	Uom3id      int           `gorm:"column:uom3id" json:"uom3id"`
	Grqty       float64       `gorm:"column:grqty" json:"grqty"`
	Grqty2      float64       `gorm:"column:grqty2" json:"grqty2"`
	Grqty3      float64       `gorm:"column:grqty3" json:"grqty3"`
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

func (Poresult) TableName() string {
	return "poresult"
}
