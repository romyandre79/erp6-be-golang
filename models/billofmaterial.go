package models

import (
  "time"
)

type Billofmaterial struct {
	Bomid int `gorm:"column:bomid;primaryKey" json:"bomid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Bomversion string `gorm:"column:bomversion" json:"bomversion"`
	Productid int `gorm:"column:productid" json:"productid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Mesinid int `gorm:"column:mesinid" json:"mesinid"`
	Numoperator int `gorm:"column:numoperator" json:"numoperator"`
	Processprdid int `gorm:"column:processprdid" json:"processprdid"`
	Bomdate time.Time `gorm:"column:bomdate" json:"bomdate"`
	Description *string `gorm:"column:description" json:"description"`
	Sla *int `gorm:"column:sla" json:"sla"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Billofmaterial) TableName() string {
	return "billofmaterial"
}
