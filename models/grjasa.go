package models

import (
  "time"
)

type Grjasa struct {
	Grjasaid int `gorm:"column:grjasaid;primaryKey" json:"grjasaid"`
	Grheaderid int `gorm:"column:grheaderid" json:"grheaderid"`
	Pojasaid int `gorm:"column:pojasaid" json:"pojasaid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Reqdate time.Time `gorm:"column:reqdate" json:"reqdate"`
	Sloctoid int `gorm:"column:sloctoid" json:"sloctoid"`
	Mesinid *int `gorm:"column:mesinid" json:"mesinid"`
	Description *string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Pojasa Pojasa `gorm:"foreignKey:pojasaid;references:pojasaid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Sloc Sloc `gorm:"foreignKey:sloctoid;references:slocid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Mesin Mesin `gorm:"foreignKey:mesinid;references:mesinid"`
}

func (Grjasa) TableName() string {
	return "grjasa"
}
