package models

import (
	"time"
)

type Expeditionapgr struct {
	Expeditionapgrid int           `gorm:"column:expeditionapgrid;primaryKey" json:"expeditionapgrid"`
	Expeditionapid   int           `gorm:"column:expeditionapid" json:"expeditionapid"`
	Grheaderid       int           `gorm:"column:grheaderid" json:"grheaderid"`
	Grdetailid       int           `gorm:"column:grdetailid" json:"grdetailid"`
	Productid        int           `gorm:"column:productid" json:"productid"`
	Productname      string        `gorm:"column:productname" json:"productname"`
	Qty              float64       `gorm:"column:qty" json:"qty"`
	Qty2             float64       `gorm:"column:qty2" json:"qty2"`
	Uomid            int           `gorm:"column:uomid" json:"uomid"`
	Uom2id           int           `gorm:"column:uom2id" json:"uom2id"`
	Nilaibeban       float64       `gorm:"column:nilaibeban" json:"nilaibeban"`
	Updatedate       time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product          Product       `gorm:"foreignKey:productid;references:productid"`
	Uom              Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2             Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Expeditionapgr) TableName() string {
	return "expeditionapgr"
}
