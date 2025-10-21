package models

import (
	"time"
)

type Rfqjasa struct {
	Rfqjasaid   int       `gorm:"column:rfqjasaid" json:"rfqjasaid"`
	Rfqid       int       `gorm:"column:rfqid" json:"rfqid"`
	Prjasaid    int       `gorm:"column:prjasaid" json:"prjasaid"`
	Prheaderid  int       `gorm:"column:prheaderid" json:"prheaderid"`
	Productid   int       `gorm:"column:productid" json:"productid"`
	Productcode string    `gorm:"column:productcode" json:"productcode"`
	Productname string    `gorm:"column:productname" json:"productname"`
	Uomid       int       `gorm:"column:uomid" json:"uomid"`
	Qty         float64   `gorm:"column:qty" json:"qty"`
	Grqty       float64   `gorm:"column:grqty" json:"grqty"`
	Tsqty       float64   `gorm:"column:tsqty" json:"tsqty"`
	Price       float64   `gorm:"column:price" json:"price"`
	Reqdate     time.Time `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid    int       `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid     int       `gorm:"column:mesinid" json:"mesinid"`
	Description string    `gorm:"column:description" json:"description"`
}

func (Rfqjasa) TableName() string {
	return "rfqjasa"
}
