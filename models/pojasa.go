package models

import (
  "time"
)

type Pojasa struct {
	Pojasaid int `gorm:"column:pojasaid;primaryKey" json:"pojasaid"`
	Poheaderid int `gorm:"column:poheaderid" json:"poheaderid"`
	Prjasaid *int `gorm:"column:prjasaid" json:"prjasaid"`
	Prheaderid *int `gorm:"column:prheaderid" json:"prheaderid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Grqty *float64 `gorm:"column:grqty" json:"grqty"`
	Tsqty *float64 `gorm:"column:tsqty" json:"tsqty"`
	Price float64 `gorm:"column:price" json:"price"`
	Reqdate time.Time `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid int `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid *int `gorm:"column:mesinid" json:"mesinid"`
	Tolerance float64 `gorm:"column:tolerance" json:"tolerance"`
	Description *string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Sloc Sloc `gorm:"foreignKey:sloctoid;references:slocid"`
	Mesin Mesin `gorm:"foreignKey:mesinid;references:mesinid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Pojasa) TableName() string {
	return "pojasa"
}
