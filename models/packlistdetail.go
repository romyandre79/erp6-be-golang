package models

import (
	"time"
)

type Packlistdetail struct {
	Packlistdetailid int           `gorm:"column:packlistdetailid;primaryKey" json:"packlistdetailid"`
	Packinglistid    int           `gorm:"column:packinglistid" json:"packinglistid"`
	Certoaid         *int          `gorm:"column:certoaid" json:"certoaid"`
	Productid        int           `gorm:"column:productid" json:"productid"`
	Productname      *string       `gorm:"column:productname" json:"productname"`
	Productcode      *string       `gorm:"column:productcode" json:"productcode"`
	Certoano         *int          `gorm:"column:certoano" json:"certoano"`
	Batchno          *string       `gorm:"column:batchno" json:"batchno"`
	Lotno            *string       `gorm:"column:lotno" json:"lotno"`
	Qty              float64       `gorm:"column:qty" json:"qty"`
	Uomid            int           `gorm:"column:uomid" json:"uomid"`
	Uomcode          *string       `gorm:"column:uomcode" json:"uomcode"`
	Qty2             float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id           int           `gorm:"column:uom2id" json:"uom2id"`
	Gross            *float64      `gorm:"column:gross" json:"gross"`
	Grossuomid       *int          `gorm:"column:grossuom" json:"grossuom"`
	Net              *float64      `gorm:"column:net" json:"net"`
	Netuomid         *int          `gorm:"column:netuom" json:"netuom"`
	Gidetailid       *int          `gorm:"column:gidetailid" json:"gidetailid"`
	Sodetailid       *int          `gorm:"column:sodetailid" json:"sodetailid"`
	Description      *string       `gorm:"column:description" json:"description"`
	Updatedate       time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Certoa           Certoa        `gorm:"foreignKey:certoaid;references:certoaid"`
	Grossuom         Unitofmeasure `gorm:"foreignKey:grossuom;references:unitofmeasureid"`
	Netuom           Unitofmeasure `gorm:"foreignKey:netuom;references:unitofmeasureid"`
	Uom              Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2             Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
}

func (Packlistdetail) TableName() string {
	return "packlistdetail"
}
