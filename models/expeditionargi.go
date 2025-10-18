package models

import (
	"time"
)

type Expeditionargi struct {
	Expeditionargiid int           `gorm:"column:expeditionargiid;primaryKey" json:"expeditionargiid"`
	Expeditionarid   int           `gorm:"column:expeditionarid" json:"expeditionarid"`
	Giheaderid       int           `gorm:"column:giheaderid" json:"giheaderid"`
	Gidetailid       int           `gorm:"column:gidetailid" json:"gidetailid"`
	Productid        int           `gorm:"column:productid" json:"productid"`
	Productname      string        `gorm:"column:productname" json:"productname"`
	Qty              float64       `gorm:"column:qty" json:"qty"`
	Qty2             float64       `gorm:"column:qty2" json:"qty2"`
	Uomid            int           `gorm:"column:uomid" json:"uomid"`
	Uom2id           int           `gorm:"column:uom2id" json:"uom2id"`
	Nilaibeban       float64       `gorm:"column:nilaibeban" json:"nilaibeban"`
	Updatedate       time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Giheader         Giheader      `gorm:"foreignKey:giheaderid;references:giheaderid"`
	Product          Product       `gorm:"foreignKey:productid;references:productid"`
	Uom              Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2             Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Expeditionargi) TableName() string {
	return "expeditionargi"
}
