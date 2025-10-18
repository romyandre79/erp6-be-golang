package models

import (
  "time"
)

type Productplandetail struct {
	Productplandetailid int `gorm:"column:productplandetailid;primaryKey" json:"productplandetailid"`
	Productplanid int `gorm:"column:productplanid" json:"productplanid"`
	Productplanfgid *int `gorm:"column:productplanfgid" json:"productplanfgid"`
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
	Sloctoid int `gorm:"column:sloctoid" json:"sloctoid"`
	Slocfromid int `gorm:"column:slocfromid" json:"slocfromid"`
	Qtyres *float64 `gorm:"column:qtyres" json:"qtyres"`
	Qtyres2 *float64 `gorm:"column:qtyres2" json:"qtyres2"`
	Qtyres3 *float64 `gorm:"column:qtyres3" json:"qtyres3"`
	Qtyfr *float64 `gorm:"column:qtyfr" json:"qtyfr"`
	Qtyfr2 *float64 `gorm:"column:qtyfr2" json:"qtyfr2"`
	Qtyfr3 *float64 `gorm:"column:qtyfr3" json:"qtyfr3"`
	Qtytrf *float64 `gorm:"column:qtytrf" json:"qtytrf"`
	Qtytrf2 *float64 `gorm:"column:qtytrf2" json:"qtytrf2"`
	Qtytrf3 *float64 `gorm:"column:qtytrf3" json:"qtytrf3"`
	Isread *int `gorm:"column:isread" json:"isread"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Billofmaterial Billofmaterial `gorm:"foreignKey:bomid;references:bomid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
}

func (Productplandetail) TableName() string {
	return "productplandetail"
}
