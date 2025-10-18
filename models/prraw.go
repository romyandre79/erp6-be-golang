package models

import (
	"time"
)

type Prraw struct {
	Prrawid          int           `gorm:"column:prrawid;primaryKey" json:"prrawid"`
	Prheaderid       int           `gorm:"column:prheaderid" json:"prheaderid"`
	Formrequestrawid int           `gorm:"column:formrequestrawid" json:"formrequestrawid"`
	Productid        int           `gorm:"column:productid" json:"productid"`
	Qty              float64       `gorm:"column:qty" json:"qty"`
	Productcode      *string       `gorm:"column:productcode" json:"productcode"`
	Productname      *string       `gorm:"column:productname" json:"productname"`
	Uomid            int           `gorm:"column:uomid" json:"uomid"`
	Qty2             float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id           int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3             *float64      `gorm:"column:qty3" json:"qty3"`
	Uom3id           *int          `gorm:"column:uom3id" json:"uom3id"`
	Reqdate          time.Time     `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid         int           `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid          *int          `gorm:"column:mesinid" json:"mesinid"`
	Poqty            *float64      `gorm:"column:poqty" json:"poqty"`
	Poqty2           *float64      `gorm:"column:poqty2" json:"poqty2"`
	Poqty3           *float64      `gorm:"column:poqty3" json:"poqty3"`
	Description      *string       `gorm:"column:description" json:"description"`
	Mesin            Mesin         `gorm:"foreignKey:mesinid;references:mesinid"`
	Product          Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc             Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
	Uom              Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2             Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3             Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Prraw) TableName() string {
	return "prraw"
}
