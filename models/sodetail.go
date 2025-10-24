package models

import (
	"time"
)

type Sodetail struct {
	Sodetailid     int            `gorm:"column:sodetailid;primaryKey" json:"sodetailid"`
	Soheaderid     int            `gorm:"column:soheaderid" json:"soheaderid"`
	Productid      int            `gorm:"column:productid" json:"productid"`
	Delvdate       time.Time      `gorm:"column:delvdate" json:"delvdate"`
	Qty            float64        `gorm:"column:qty" json:"qty"`
	Uomid          int            `gorm:"column:uomid" json:"uomid"`
	Qty2           float64        `gorm:"column:qty2" json:"qty2"`
	Uom2id         int            `gorm:"column:uom2id" json:"uom2id"`
	Price          float64        `gorm:"column:price" json:"price"`
	Qty3           float64        `gorm:"column:qty3" json:"qty3"`
	Uom3id         float64        `gorm:"column:uom3id" json:"uom3id"`
	Bomid          int            `gorm:"column:bomid" json:"bomid"`
	Toleransi      float64        `gorm:"column:toleransi" json:"toleransi"`
	Slocid         int            `gorm:"column:slocid" json:"slocid"`
	Itemnote       string         `gorm:"column:itemnote" json:"itemnote"`
	Giqty          float64        `gorm:"column:giqty" json:"giqty"`
	Giqty2         float64        `gorm:"column:giqty2" json:"giqty2"`
	Invqty         float64        `gorm:"column:invqty" json:"invqty"`
	Invqty2        float64        `gorm:"column:invqty2" json:"invqty2"`
	Sppqty         float64        `gorm:"column:sppqty" json:"sppqty"`
	Sppqty2        float64        `gorm:"column:sppqty2" json:"sppqty2"`
	Opqty          float64        `gorm:"column:opqty" json:"opqty"`
	Opqty2         float64        `gorm:"column:opqty2" json:"opqty2"`
	Vrqty          float64        `gorm:"column:vrqty" json:"vrqty"`
	Vrqty2         float64        `gorm:"column:vrqty2" json:"vrqty2"`
	Billofmaterial Billofmaterial `gorm:"foreignKey:bomid;references:bomid"`
	Product        Product        `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure  Unitofmeasure  `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Sodetail) TableName() string {
	return "sodetail"
}
