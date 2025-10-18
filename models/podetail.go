package models

import (
	"time"
)

type Podetail struct {
	Podetailid    int           `gorm:"column:podetailid;primaryKey" json:"podetailid"`
	Poheaderid    int           `gorm:"column:poheaderid" json:"poheaderid"`
	Prrawid       *int          `gorm:"column:prrawid" json:"prrawid"`
	Prheaderid    *int          `gorm:"column:prheaderid" json:"prheaderid"`
	Productid     int           `gorm:"column:productid" json:"productid"`
	Productcode   *string       `gorm:"column:productcode" json:"productcode"`
	Productname   *string       `gorm:"column:productname" json:"productname"`
	Qty           float64       `gorm:"column:qty" json:"qty"`
	Qty2          float64       `gorm:"column:qty2" json:"qty2"`
	Qty3          *float64      `gorm:"column:qty3" json:"qty3"`
	Grqty         float64       `gorm:"column:grqty" json:"grqty"`
	Grqty2        *float64      `gorm:"column:grqty2" json:"grqty2"`
	Grqty3        *float64      `gorm:"column:grqty3" json:"grqty3"`
	Tsqty         float64       `gorm:"column:tsqty" json:"tsqty"`
	Tsqty2        float64       `gorm:"column:tsqty2" json:"tsqty2"`
	Tsqty3        *float64      `gorm:"column:tsqty3" json:"tsqty3"`
	Uomid         int           `gorm:"column:uomid" json:"uomid"`
	Uom2id        int           `gorm:"column:uom2id" json:"uom2id"`
	Uom3id        *int          `gorm:"column:uom3id" json:"uom3id"`
	Price         float64       `gorm:"column:price" json:"price"`
	Arrivedate    *time.Time    `gorm:"column:arrivedate" json:"arrivedate"`
	Toleransiup   float64       `gorm:"column:toleransiup" json:"toleransiup"`
	Toleransidown float64       `gorm:"column:toleransidown" json:"toleransidown"`
	Slocid        int           `gorm:"column:slocid" json:"slocid"`
	Sloctoid      int           `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid       *int          `gorm:"column:mesinid" json:"mesinid"`
	Requestedbyid *int          `gorm:"column:requestedbyid" json:"requestedbyid"`
	Itemnote      *string       `gorm:"column:itemnote" json:"itemnote"`
	Updatedate    time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Prheader      Prheader      `gorm:"foreignKey:prheaderid;references:prheaderid"`
	Product       Product       `gorm:"foreignKey:productid;references:productid"`
	Requestedby   Requestedby   `gorm:"foreignKey:requestedbyid;references:requestedbyid"`
	Sloc          Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Uom           Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2          Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Podetail) TableName() string {
	return "podetail"
}
