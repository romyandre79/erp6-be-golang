package models

import (
  "time"
)

type Productplant struct {
	Productplantid int `gorm:"column:productplantid;primaryKey" json:"productplantid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Issource int8 `gorm:"column:issource" json:"issource"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Sloc Sloc `gorm:"foreignKey:slocid;references:slocid"`
}

func (Productplant) TableName() string {
	return "productplant"
}
