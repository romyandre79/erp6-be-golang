package models

import (
	"time"
)

type Rawschedule struct {
	Rawscheduleid int       `gorm:"column:rawscheduleid;primaryKey" json:"rawscheduleid"`
	Poheaderid    int       `gorm:"column:poheaderid" json:"poheaderid"`
	Podetailid    int       `gorm:"column:podetailid" json:"podetailid"`
	Addressbookid int       `gorm:"column:addressbookid" json:"addressbookid"`
	Docno         string    `gorm:"column:docno" json:"docno"`
	Slocid        int       `gorm:"column:slocid" json:"slocid"`
	Productid     int       `gorm:"column:productid" json:"productid"`
	Qty           float64   `gorm:"column:qty" json:"qty"`
	Qty2          float64   `gorm:"column:qty2" json:"qty2"`
	Qty3          float64   `gorm:"column:qty3" json:"qty3"`
	Uomid         int       `gorm:"column:uomid" json:"uomid"`
	Uom2id        int       `gorm:"column:uom2id" json:"uom2id"`
	Uom3id        int       `gorm:"column:uom3id" json:"uom3id"`
	Arrivedate    time.Time `gorm:"column:arrivedate" json:"arrivedate"`
	Mesinid       int       `gorm:"column:mesinid" json:"mesinid"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Rawschedule) TableName() string {
	return "rawschedule"
}
