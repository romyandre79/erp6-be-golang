package models

import (
	"time"
)

type Productoutputwaste struct {
	Productoutputwasteid int           `gorm:"column:productoutputwasteid;primaryKey" json:"productoutputwasteid"`
	Productoutputid      int           `gorm:"column:productoutputid" json:"productoutputid"`
	Productoutputfgid    int           `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Productid            int           `gorm:"column:productid" json:"productid"`
	Productname          string        `gorm:"column:productname" json:"productname"`
	Qty                  float64       `gorm:"column:qty" json:"qty"`
	Uomid                int           `gorm:"column:uomid" json:"uomid"`
	Qty2                 float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id               int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3                 float64       `gorm:"column:qty3" json:"qty3"`
	Uom3id               int           `gorm:"column:uom3id" json:"uom3id"`
	Slocid               int           `gorm:"column:slocid" json:"slocid"`
	Storagebinid         int           `gorm:"column:storagebinid" json:"storagebinid"`
	Itemnote             string        `gorm:"column:itemnote" json:"itemnote"`
	Updatedate           time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product              Product       `gorm:"foreignKey:productid;references:productid"`
	Uom                  Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2                 Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Productoutputwaste) TableName() string {
	return "productoutputwaste"
}
