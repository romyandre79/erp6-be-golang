package models

import (
	"time"
)

type Gireturdetail struct {
	Gireturdetailid int           `gorm:"column:gireturdetailid;primaryKey" json:"gireturdetailid"`
	Gireturid       int           `gorm:"column:gireturid" json:"gireturid"`
	Productid       int           `gorm:"column:productid" json:"productid"`
	Qty             float64       `gorm:"column:qty" json:"qty"`
	Qty2            float64       `gorm:"column:qty2" json:"qty2"`
	Qty3            float64       `gorm:"column:qty3" json:"qty3"`
	Qty4            float64       `gorm:"column:qty4" json:"qty4"`
	Uomid           int           `gorm:"column:uomid" json:"uomid"`
	Uom2id          int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id          *int          `gorm:"column:uom3id" json:"uom3id"`
	Uom4id          *int          `gorm:"column:uom4id" json:"uom4id"`
	Gidetailid      int           `gorm:"column:gidetailid" json:"gidetailid"`
	Slocid          int           `gorm:"column:slocid" json:"slocid"`
	Storagebinid    int           `gorm:"column:storagebinid" json:"storagebinid"`
	Itemnote        *string       `gorm:"column:itemnote" json:"itemnote"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Gidetail        Gidetail      `gorm:"foreignKey:gidetailid;references:gidetailid"`
	Product         Product       `gorm:"foreignKey:productid;references:productid"`
	Storagebin      Storagebin    `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Sloc            Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom             Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2            Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3            Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
	Uom4            Unitofmeasure `gorm:"foreignKey:uom4id;references:unitofmeasureid"`
}

func (Gireturdetail) TableName() string {
	return "gireturdetail"
}
