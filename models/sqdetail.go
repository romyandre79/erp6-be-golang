package models

import (
	"time"
)

type Sqdetail struct {
	Sqdetailid int           `gorm:"column:sqdetailid;primaryKey" json:"sqdetailid"`
	Sqheaderid int           `gorm:"column:sqheaderid" json:"sqheaderid"`
	Productid  int           `gorm:"column:productid" json:"productid"`
	Delvdate   time.Time     `gorm:"column:delvdate" json:"delvdate"`
	Qty        float64       `gorm:"column:qty" json:"qty"`
	Uomid      int           `gorm:"column:uomid" json:"uomid"`
	Qty2       float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id     int           `gorm:"column:uom2id" json:"uom2id"`
	Price      float64       `gorm:"column:price" json:"price"`
	Qty3       *float64      `gorm:"column:qty3" json:"qty3"`
	Uom3id     *int          `gorm:"column:uom3id" json:"uom3id"`
	Slocid     *int          `gorm:"column:slocid" json:"slocid"`
	Toleransi  float64       `gorm:"column:toleransi" json:"toleransi"`
	Itemnote   string        `gorm:"column:itemnote" json:"itemnote"`
	Updatedate time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product    Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc       Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom        Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2       Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3       Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Sqdetail) TableName() string {
	return "sqdetail"
}
