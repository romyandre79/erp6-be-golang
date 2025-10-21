package models

import (
	"time"
)

type Formrequestjasa struct {
	Formrequestjasaid int           `gorm:"column:formrequestjasaid;primaryKey" json:"formrequestjasaid"`
	Formrequestid     int           `gorm:"column:formrequestid" json:"formrequestid"`
	Productid         int           `gorm:"column:productid" json:"productid"`
	Productcode       string        `gorm:"column:productcode" json:"productcode"`
	Productname       string        `gorm:"column:productname" json:"productname"`
	Qty               float64       `gorm:"column:qty" json:"qty"`
	Uomid             int           `gorm:"column:uomid" json:"uomid"`
	Mesinid           int           `gorm:"column:mesinid" json:"mesinid"`
	Reqdate           time.Time     `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid          int           `gorm:"column:sloctoid" json:"sloctoid"`
	Prqty             float64       `gorm:"column:prqty" json:"prqty"`
	Tsqty             float64       `gorm:"column:tsqty" json:"tsqty"`
	Description       string        `gorm:"column:description" json:"description"`
	Updatedate        time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Mesin             Mesin         `gorm:"foreignKey:mesinid;references:mesinid"`
	Product           Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc              Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
	Unitofmeasure     Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Formrequestjasa) TableName() string {
	return "formrequestjasa"
}
