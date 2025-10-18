package models

import (
	"time"
)

type Notagrrdetail struct {
	Notagrrdetailid int           `gorm:"column:notagrrdetailid;primaryKey" json:"notagrrdetailid"`
	Notagrrid       int           `gorm:"column:notagrrid" json:"notagrrid"`
	Notagrrgrrid    *int          `gorm:"column:notagrrgrrid" json:"notagrrgrrid"`
	Productid       int           `gorm:"column:productid" json:"productid"`
	Productname     *string       `gorm:"column:productname" json:"productname"`
	Qty             float64       `gorm:"column:qty" json:"qty"`
	Qty2            float64       `gorm:"column:qty2" json:"qty2"`
	Qty3            float64       `gorm:"column:qty3" json:"qty3"`
	Qty4            float64       `gorm:"column:qty4" json:"qty4"`
	Uomid           int           `gorm:"column:uomid" json:"uomid"`
	Slocid          int           `gorm:"column:slocid" json:"slocid"`
	Uom2id          int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id          *int          `gorm:"column:uom3id" json:"uom3id"`
	Uom4id          *int          `gorm:"column:uom4id" json:"uom4id"`
	Price           float64       `gorm:"column:price" json:"price"`
	Currencyid      int           `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue       float64       `gorm:"column:ratevalue" json:"ratevalue"`
	Dpp             float64       `gorm:"column:dpp" json:"dpp"`
	Total           *float64      `gorm:"column:total" json:"total"`
	Detailnote      *string       `gorm:"column:detailnote" json:"detailnote"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Currency        Currency      `gorm:"foreignKey:currencyid;references:currencyid"`
	Product         Product       `gorm:"foreignKey:productid;references:productid"`
	Uom             Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2            Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3            Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
	Uom4            Unitofmeasure `gorm:"foreignKey:uom4id;references:unitofmeasureid"`
}

func (Notagrrdetail) TableName() string {
	return "notagrrdetail"
}
