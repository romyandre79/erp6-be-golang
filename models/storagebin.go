package models

import (
  "time"
)

type Storagebin struct {
	Storagebinid int `gorm:"column:storagebinid;primaryKey" json:"storagebinid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Description string `gorm:"column:description" json:"description"`
	Ismultiproduct int8 `gorm:"column:ismultiproduct" json:"ismultiproduct"`
	Qtymax float64 `gorm:"column:qtymax" json:"qtymax"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Storagebin) TableName() string {
	return "storagebin"
}
