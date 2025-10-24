package models

import (
	"time"
)

type Product struct {
	Productid       int           `gorm:"column:productid;primaryKey" json:"productid"`
	Companyid       int           `gorm:"column:companyid" json:"companyid"`
	Productcode     string        `gorm:"column:productcode" json:"productcode"`
	Productname     string        `gorm:"column:productname" json:"productname"`
	Productdesc     string        `gorm:"column:productdesc" json:"productdesc"`
	Productpic      string        `gorm:"column:productpic" json:"productpic"`
	Barcode         string        `gorm:"column:barcode" json:"barcode"`
	Isstock         int8          `gorm:"column:isstock" json:"isstock"`
	Isasset         int8          `gorm:"column:isasset" json:"isasset"`
	Materialtypeid  int           `gorm:"column:materialtypeid" json:"materialtypeid"`
	Thickness       float64       `gorm:"column:thickness" json:"thickness"`
	Width           float64       `gorm:"column:width" json:"width"`
	Length          float64       `gorm:"column:length" json:"length"`
	Colourid        int           `gorm:"column:colourid" json:"colourid"`
	Qty1            float64       `gorm:"column:qty1" json:"qty1"`
	Uom1            int           `gorm:"column:uom1" json:"uom1"`
	Qty2            float64       `gorm:"column:qty2" json:"qty2"`
	Uom2            int           `gorm:"column:uom2" json:"uom2"`
	Qty3            float64       `gorm:"column:qty3" json:"qty3"`
	Uom3            int           `gorm:"column:uom3" json:"uom3"`
	Materialgroupid int           `gorm:"column:materialgroupid" json:"materialgroupid"`
	Sled            int           `gorm:"column:sled" json:"sled"`
	Isautolot       int8          `gorm:"column:isautolot" json:"isautolot"`
	Recordstatus    int8          `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Company         Company       `gorm:"foreignKey:companyid;references:companyid"`
	Materialgroup   Materialgroup `gorm:"foreignKey:materialgroupid;references:materialgroupid"`
	Materialtype    Materialtype  `gorm:"foreignKey:materialtypeid;references:materialtypeid"`
}

func (Product) TableName() string {
	return "product"
}
