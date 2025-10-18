package models

import (
  "time"
)

type Pricecategory struct {
	Pricecategoryid int `gorm:"column:pricecategoryid;primaryKey" json:"pricecategoryid"`
	Categoryname string `gorm:"column:categoryname" json:"categoryname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Pricecategory) TableName() string {
	return "pricecategory"
}
