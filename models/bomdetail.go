package models

import (
  "time"
)

type Bomdetail struct {
	Bomdetailid int `gorm:"column:bomdetailid;primaryKey" json:"bomdetailid"`
	Bomid int `gorm:"column:bomid" json:"bomid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Productbomid *int `gorm:"column:productbomid" json:"productbomid"`
	Productparentid *int `gorm:"column:productparentid" json:"productparentid"`
	Description *string `gorm:"column:description" json:"description"`
	Sla *int `gorm:"column:sla" json:"sla"`
	Updatedate *time.Time `gorm:"column:updatedate" json:"updatedate"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Bomdetail) TableName() string {
	return "bomdetail"
}
