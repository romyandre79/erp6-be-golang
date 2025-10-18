package models

import (
  "time"
)

type Stdfohul struct {
	Stdfohulid int `gorm:"column:stdfohulid;primaryKey" json:"stdfohulid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Valued float64 `gorm:"column:valued" json:"valued"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Stdfohul) TableName() string {
	return "stdfohul"
}
