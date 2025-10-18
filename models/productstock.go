package models

import (
	"time"
)

type Productstock struct {
	Productstockid     int           `gorm:"column:productstockid;primaryKey" json:"productstockid"`
	Productid          int           `gorm:"column:productid" json:"productid"`
	Addressbookid      *int          `gorm:"column:addressbookid" json:"addressbookid"`
	Productcode        *string       `gorm:"column:productcode" json:"productcode"`
	Productname        string        `gorm:"column:productname" json:"productname"`
	Slocid             int           `gorm:"column:slocid" json:"slocid"`
	Sloccode           string        `gorm:"column:sloccode" json:"sloccode"`
	Storagebinid       int           `gorm:"column:storagebinid" json:"storagebinid"`
	Storagedesc        string        `gorm:"column:storagedesc" json:"storagedesc"`
	Qty                *float64      `gorm:"column:qty" json:"qty"`
	Qty2               *float64      `gorm:"column:qty2" json:"qty2"`
	Qty3               *float64      `gorm:"column:qty3" json:"qty3"`
	Uomid              int           `gorm:"column:uomid" json:"uomid"`
	Uomcode            string        `gorm:"column:uomcode" json:"uomcode"`
	Uom2id             int           `gorm:"column:uom2id" json:"uom2id"`
	Uom2code           string        `gorm:"column:uom2code" json:"uom2code"`
	Uom3id             *int          `gorm:"column:uom3id" json:"uom3id"`
	Uom3code           *string       `gorm:"column:uom3code" json:"uom3code"`
	Materialstatusid   *int          `gorm:"column:materialstatusid" json:"materialstatusid"`
	Materialstatusname *string       `gorm:"column:materialstatusname" json:"materialstatusname"`
	Lotno              *string       `gorm:"column:lotno" json:"lotno"`
	Qtyalokasi         float64       `gorm:"column:qtyalokasi" json:"qtyalokasi"`
	Qtyinprogress      float64       `gorm:"column:qtyinprogress" json:"qtyinprogress"`
	Updatedate         time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product            Product       `gorm:"foreignKey:productid;references:productid"`
	Storagebin         Storagebin    `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Sloc               Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom                Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2               Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Productstock) TableName() string {
	return "productstock"
}
