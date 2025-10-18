package models

import (

)

type Sfdetail struct {
	Sfdetailid int `gorm:"column:sfdetailid;primaryKey" json:"sfdetailid"`
	Sfheaderid int `gorm:"column:sfheaderid" json:"sfheaderid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Price float64 `gorm:"column:price" json:"price"`
	Itemnote *string `gorm:"column:itemnote" json:"itemnote"`
	Vrqty *float64 `gorm:"column:vrqty" json:"vrqty"`
	Vrqty2 *float64 `gorm:"column:vrqty2" json:"vrqty2"`
}

func (Sfdetail) TableName() string {
	return "sfdetail"
}
