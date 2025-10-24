package models

import (
	"time"
)

type Gidetail struct {
	Gidetailid        int           `gorm:"column:gidetailid;primaryKey" json:"gidetailid"`
	Giheaderid        int           `gorm:"column:giheaderid" json:"giheaderid"`
	Productid         int           `gorm:"column:productid" json:"productid"`
	Productcode       string        `gorm:"column:productcode" json:"productcode"`
	Productname       string        `gorm:"column:productname" json:"productname"`
	Qty               float64       `gorm:"column:qty" json:"qty"`
	Uomid             int           `gorm:"column:uomid" json:"uomid"`
	Qty2              float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id            int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3              float64       `gorm:"column:qty3" json:"qty3"`
	Uom3id            int           `gorm:"column:uom3id" json:"uom3id"`
	Qty4              float64       `gorm:"column:qty4" json:"qty4"`
	Uom4id            int           `gorm:"column:uom4id" json:"uom4id"`
	Slocid            int           `gorm:"column:slocid" json:"slocid"`
	Storagebinid      int           `gorm:"column:storagebinid" json:"storagebinid"`
	Lotno             string        `gorm:"column:lotno" json:"lotno"`
	Certoaid          int           `gorm:"column:certoaid" json:"certoaid"`
	Sodetailid        int           `gorm:"column:sodetailid" json:"sodetailid"`
	Qtyinv            float64       `gorm:"column:qtyinv" json:"qtyinv"`
	Qtyinv2           float64       `gorm:"column:qtyinv2" json:"qtyinv2"`
	Qtyinv3           float64       `gorm:"column:qtyinv3" json:"qtyinv3"`
	Qtyinv4           float64       `gorm:"column:qtyinv4" json:"qtyinv4"`
	Productoutputid   int           `gorm:"column:productoutputid" json:"productoutputid"`
	Productstockdetid int           `gorm:"column:productstockdetid" json:"productstockdetid"`
	Qtyretur          float64       `gorm:"column:qtyretur" json:"qtyretur"`
	Itemnote          string        `gorm:"column:itemnote" json:"itemnote"`
	Updatedate        time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Product           Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc              Sloc          `gorm:"foreignKey:slocid;references:slocid"`
	Sodetail          Sodetail      `gorm:"foreignKey:sodetailid;references:sodetailid"`
	Storagebin        Storagebin    `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Uom               Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2              Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3              Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
	Uom4              Unitofmeasure `gorm:"foreignKey:uom4id;references:unitofmeasureid"`
}

func (Gidetail) TableName() string {
	return "gidetail"
}
