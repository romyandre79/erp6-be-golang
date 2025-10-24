package models

import (
	"time"
)

type Deliveryschedule struct {
	Deliveryscheduleid int           `gorm:"column:deliveryscheduleid;primaryKey" json:"deliveryscheduleid"`
	Soheaderid         int           `gorm:"column:soheaderid" json:"soheaderid"`
	Sodetailid         int           `gorm:"column:sodetailid" json:"sodetailid"`
	Addressbookid      int           `gorm:"column:addressbookid" json:"addressbookid"`
	Docno              string        `gorm:"column:docno" json:"docno"`
	Addressid          int           `gorm:"column:addressid" json:"addressid"`
	Delvdate           time.Time     `gorm:"column:delvdate" json:"delvdate"`
	Productid          int           `gorm:"column:productid" json:"productid"`
	Qty                float64       `gorm:"column:qty" json:"qty"`
	Qty2               float64       `gorm:"column:qty2" json:"qty2"`
	Qty3               float64       `gorm:"column:qty3" json:"qty3"`
	Qty4               float64       `gorm:"column:qty4" json:"qty4"`
	Uomid              int           `gorm:"column:uomid" json:"uomid"`
	Uom2id             int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id             int           `gorm:"column:uom3id" json:"uom3id"`
	Uom4id             int           `gorm:"column:uom4id" json:"uom4id"`
	Updatedate         time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Addressbook        Addressbook   `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Address            Address       `gorm:"foreignKey:addressid;references:addressid"`
	Product            Product       `gorm:"foreignKey:productid;references:productid"`
	Sodetail           Sodetail      `gorm:"foreignKey:sodetailid;references:sodetailid"`
	Soheader           Soheader      `gorm:"foreignKey:soheaderid;references:soheaderid"`
	Uom                Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2               Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Deliveryschedule) TableName() string {
	return "deliveryschedule"
}
