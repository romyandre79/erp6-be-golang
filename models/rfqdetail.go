package models

import (
	"time"
)

type Rfqdetail struct {
	Rfqdetailid   int       `gorm:"column:rfqdetailid;primaryKey" json:"rfqdetailid"`
	Rfqid         int       `gorm:"column:rfqid" json:"rfqid"`
	Prrawid       int       `gorm:"column:prrawid" json:"prrawid"`
	Prheaderid    int       `gorm:"column:prheaderid" json:"prheaderid"`
	Productid     int       `gorm:"column:productid" json:"productid"`
	Qty           float64   `gorm:"column:qty" json:"qty"`
	Qty2          float64   `gorm:"column:qty2" json:"qty2"`
	Grqty         float64   `gorm:"column:grqty" json:"grqty"`
	Grqty2        float64   `gorm:"column:grqty2" json:"grqty2"`
	Tsqty         float64   `gorm:"column:tsqty" json:"tsqty"`
	Tsqty2        float64   `gorm:"column:tsqty2" json:"tsqty2"`
	Qtyres        float64   `gorm:"column:qtyres" json:"qtyres"`
	Uomid         int       `gorm:"column:uomid" json:"uomid"`
	Uom2id        int       `gorm:"column:uom2id" json:"uom2id"`
	Price         float64   `gorm:"column:price" json:"price"`
	Arrivedate    time.Time `gorm:"column:arrivedate" json:"arrivedate"`
	Toleransiup   float64   `gorm:"column:toleransiup" json:"toleransiup"`
	Toleransidown float64   `gorm:"column:toleransidown" json:"toleransidown"`
	Slocid        int       `gorm:"column:slocid" json:"slocid"`
	Requestedbyid int       `gorm:"column:requestedbyid" json:"requestedbyid"`
	Itemnote      string    `gorm:"column:itemnote" json:"itemnote"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Rfqdetail) TableName() string {
	return "rfqdetail"
}
