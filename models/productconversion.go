package models

import (
	"time"
)

type Productconversion struct {
	Productconversionid int           `gorm:"column:productconversionid;primaryKey" json:"productconversionid"`
	Productid           int           `gorm:"column:productid" json:"productid"`
	Qty                 float64       `gorm:"column:qty" json:"qty"`
	Qty2                float64       `gorm:"column:qty2" json:"qty2"`
	Qty3                *float64      `gorm:"column:qty3" json:"qty3"`
	Uomid               int           `gorm:"column:uomid" json:"uomid"`
	Uom2id              int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id              *int          `gorm:"column:uom3id" json:"uom3id"`
	Updatedatae         time.Time     `gorm:"column:updatedatae" json:"updatedatae"`
	Uom                 Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2                Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3                Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Productconversion) TableName() string {
	return "productconversion"
}
