package models

import (
	"time"
)

type Bsdetail struct {
	Bsdetailid       int            `gorm:"column:bsdetailid;primaryKey" json:"bsdetailid"`
	Bsheaderid       int            `gorm:"column:bsheaderid" json:"bsheaderid"`
	Productid        int            `gorm:"column:productid" json:"productid"`
	Productcode      string         `gorm:"column:productcode" json:"productcode"`
	Productname      string         `gorm:"column:productname" json:"productname"`
	Uomid            int            `gorm:"column:uomid" json:"uomid"`
	Uom2id           int            `gorm:"column:uom2id" json:"uom2id"`
	Uom3id           int            `gorm:"column:uom3id" json:"uom3id"`
	Uom4id           int            `gorm:"column:uom4id" json:"uom4id"`
	Qty              float64        `gorm:"column:qty" json:"qty"`
	Qty2             float64        `gorm:"column:qty2" json:"qty2"`
	Qty3             float64        `gorm:"column:qty3" json:"qty3"`
	Qty4             float64        `gorm:"column:qty4" json:"qty4"`
	Storagebinid     int            `gorm:"column:storagebinid" json:"storagebinid"`
	Materialstatusid int            `gorm:"column:materialstatusid" json:"materialstatusid"`
	Ownershipid      int            `gorm:"column:ownershipid" json:"ownershipid"`
	Lotno            string         `gorm:"column:lotno" json:"lotno"`
	Buyprice         float64        `gorm:"column:buyprice" json:"buyprice"`
	Buydate          time.Time      `gorm:"column:buydate" json:"buydate"`
	Currencyid       int            `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate     float64        `gorm:"column:currencyrate" json:"currencyrate"`
	Itemnote         string         `gorm:"column:itemnote" json:"itemnote"`
	Updatedate       time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Materialstatus   Materialstatus `gorm:"foreignKey:materialstatusid;references:materialstatusid"`
	Ownership        Ownership      `gorm:"foreignKey:ownershipid;references:ownershipid"`
	Product          Product        `gorm:"foreignKey:productid;references:productid"`
	Storagebin       Storagebin     `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Uom              Unitofmeasure  `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2             Unitofmeasure  `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3             Unitofmeasure  `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
	Uom4             Unitofmeasure  `gorm:"foreignKey:uom4id;references:unitofmeasureid"`
}

func (Bsdetail) TableName() string {
	return "bsdetail"
}
