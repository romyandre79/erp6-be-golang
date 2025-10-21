package models

import (
	"time"
)

type Prjasa struct {
	Prjasaid          int           `gorm:"column:prjasaid;primaryKey" json:"prjasaid"`
	Prheaderid        int           `gorm:"column:prheaderid" json:"prheaderid"`
	Formrequestjasaid int           `gorm:"column:formrequestjasaid" json:"formrequestjasaid"`
	Productid         int           `gorm:"column:productid" json:"productid"`
	Productcode       string        `gorm:"column:productcode" json:"productcode"`
	Productname       string        `gorm:"column:productname" json:"productname"`
	Qty               float64       `gorm:"column:qty" json:"qty"`
	Uomid             int           `gorm:"column:uomid" json:"uomid"`
	Mesinid           int           `gorm:"column:mesinid" json:"mesinid"`
	Reqdate           time.Time     `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid          int           `gorm:"column:sloctoid" json:"sloctoid"`
	Poqty             float64       `gorm:"column:poqty" json:"poqty"`
	Description       string        `gorm:"column:description" json:"description"`
	Updatedate        time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Mesin             Mesin         `gorm:"foreignKey:mesinid;references:mesinid"`
	Product           Product       `gorm:"foreignKey:productid;references:productid"`
	Sloc              Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
	Unitofmeasure     Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Prjasa) TableName() string {
	return "prjasa"
}
