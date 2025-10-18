package models

import (
	"time"
)

type Invoiceapdetail struct {
	Invoiceapdetailid int           `gorm:"column:invoiceapdetailid;primaryKey" json:"invoiceapdetailid"`
	Invoiceapid       int           `gorm:"column:invoiceapid" json:"invoiceapid"`
	Invoiceapgrid     *int          `gorm:"column:invoiceapgrid" json:"invoiceapgrid"`
	Slocid            int           `gorm:"column:slocid" json:"slocid"`
	Productid         int           `gorm:"column:productid" json:"productid"`
	Productname       *string       `gorm:"column:productname" json:"productname"`
	Qty               float64       `gorm:"column:qty" json:"qty"`
	Qty2              *float64      `gorm:"column:qty2" json:"qty2"`
	Qty3              *float64      `gorm:"column:qty3" json:"qty3"`
	Uomid             int           `gorm:"column:uomid" json:"uomid"`
	Uom2id            *int          `gorm:"column:uom2id" json:"uom2id"`
	Uom3id            *int          `gorm:"column:uom3id" json:"uom3id"`
	Price             *float64      `gorm:"column:price" json:"price"`
	Discount          *float64      `gorm:"column:discount" json:"discount"`
	Currencyid        int           `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue         *float64      `gorm:"column:ratevalue" json:"ratevalue"`
	Dpp               *float64      `gorm:"column:dpp" json:"dpp"`
	Total             *float64      `gorm:"column:total" json:"total"`
	Detailnote        *string       `gorm:"column:detailnote" json:"detailnote"`
	Updatedate        time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Currency          Currency      `gorm:"foreignKey:currencyid;references:currencyid"`
	Product           Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc              Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom               Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2              Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3              Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Invoiceapdetail) TableName() string {
	return "invoiceapdetail"
}
