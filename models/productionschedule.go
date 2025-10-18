package models

import (
	"time"
)

type Productionschedule struct {
	Productionscheduleid int           `gorm:"column:productionscheduleid;primaryKey" json:"productionscheduleid"`
	Productplanid        *int          `gorm:"column:productplanid" json:"productplanid"`
	Soheaderid           *int          `gorm:"column:soheaderid" json:"soheaderid"`
	Addressbookid        *int          `gorm:"column:addressbookid" json:"addressbookid"`
	Parentplanid         *int          `gorm:"column:parentplanid" json:"parentplanid"`
	Docno                string        `gorm:"column:docno" json:"docno"`
	Productplanfgid      *int          `gorm:"column:productplanfgid" json:"productplanfgid"`
	Productid            *int          `gorm:"column:productid" json:"productid"`
	Uomid                *int          `gorm:"column:uomid" json:"uomid"`
	Uom2id               *int          `gorm:"column:uom2id" json:"uom2id"`
	Uom3id               *int          `gorm:"column:uom3id" json:"uom3id"`
	Qty                  float64       `gorm:"column:qty" json:"qty"`
	Qty2                 float64       `gorm:"column:qty2" json:"qty2"`
	Qty3                 *float64      `gorm:"column:qty3" json:"qty3"`
	Processprdid         *int          `gorm:"column:processprdid" json:"processprdid"`
	Mesinid              *int          `gorm:"column:mesinid" json:"mesinid"`
	Sloctoid             *int          `gorm:"column:sloctoid" json:"sloctoid"`
	Startdate            *time.Time    `gorm:"column:startdate" json:"startdate"`
	Enddate              *time.Time    `gorm:"column:enddate" json:"enddate"`
	Updatedate           time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Addressbook          Addressbook   `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Mesin                Mesin         `gorm:"foreignKey:mesinid;references:mesinid"`
	Productplan          Productplan   `gorm:"foreignKey:productplanid;references:productplanid"`
	Productplanfg        Productplanfg `gorm:"foreignKey:productplanfgid;references:productplanfgid"`
	Processprd           Processprd    `gorm:"foreignKey:processprdid;references:processprdid"`
	Sloc                 Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
	Uom                  Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2                 Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Productionschedule) TableName() string {
	return "productionschedule"
}
