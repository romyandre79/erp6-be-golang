package models

import (
  "time"
)

type Productplanfg struct {
	Productplanfgid int `gorm:"column:productplanfgid;primaryKey" json:"productplanfgid"`
	Productplanid int `gorm:"column:productplanid" json:"productplanid"`
	Parentid *int `gorm:"column:parentid" json:"parentid"`
	Sodetailid *int `gorm:"column:sodetailid" json:"sodetailid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Bomid *int `gorm:"column:bomid" json:"bomid"`
	Processprdid *int `gorm:"column:processprdid" json:"processprdid"`
	Mesinid *int `gorm:"column:mesinid" json:"mesinid"`
	Sloctoid *int `gorm:"column:sloctoid" json:"sloctoid"`
	Qtyres *float64 `gorm:"column:qtyres" json:"qtyres"`
	Qty2res *float64 `gorm:"column:qty2res" json:"qty2res"`
	Qty3res *float64 `gorm:"column:qty3res" json:"qty3res"`
	Tsqty *float64 `gorm:"column:tsqty" json:"tsqty"`
	Ts2qty *float64 `gorm:"column:ts2qty" json:"ts2qty"`
	Ts3qty *float64 `gorm:"column:ts3qty" json:"ts3qty"`
	Isread *int8 `gorm:"column:isread" json:"isread"`
	Startdate *time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate *time.Time `gorm:"column:enddate" json:"enddate"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Productplanfg) TableName() string {
	return "productplanfg"
}
