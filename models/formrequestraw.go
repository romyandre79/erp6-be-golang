package models

import (
	"time"
)

type Formrequestraw struct {
	Formrequestrawid    int           `gorm:"column:formrequestrawid;primaryKey" json:"formrequestrawid"`
	Formrequestid       int           `gorm:"column:formrequestid" json:"formrequestid"`
	Productplandetailid int           `gorm:"column:productplandetailid" json:"productplandetailid"`
	Productid           int           `gorm:"column:productid" json:"productid"`
	Qty                 float64       `gorm:"column:qty" json:"qty"`
	Productcode         string        `gorm:"column:productcode" json:"productcode"`
	Productname         string        `gorm:"column:productname" json:"productname"`
	Uomid               int           `gorm:"column:uomid" json:"uomid"`
	Qty2                float64       `gorm:"column:qty2" json:"qty2"`
	Uom2id              int           `gorm:"column:uom2id" json:"uom2id"`
	Qty3                float64       `gorm:"column:qty3" json:"qty3"`
	Uom3id              int           `gorm:"column:uom3id" json:"uom3id"`
	Reqdate             time.Time     `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid            int           `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid             int           `gorm:"column:mesinid" json:"mesinid"`
	Prqty               float64       `gorm:"column:prqty" json:"prqty"`
	Prqty2              float64       `gorm:"column:prqty2" json:"prqty2"`
	Prqty3              float64       `gorm:"column:prqty3" json:"prqty3"`
	Tsqty               float64       `gorm:"column:tsqty" json:"tsqty"`
	Tsqty2              float64       `gorm:"column:tsqty2" json:"tsqty2"`
	Tsqty3              float64       `gorm:"column:tsqty3" json:"tsqty3"`
	Description         string        `gorm:"column:description" json:"description"`
	Updatedate          time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Mesin               Mesin         `gorm:"foreignKey:mesinid;references:mesinid"`
	Product             Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc                Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
	Uom                 Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2                Unitofmeasure `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3                Unitofmeasure `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Formrequestraw) TableName() string {
	return "formrequestraw"
}
