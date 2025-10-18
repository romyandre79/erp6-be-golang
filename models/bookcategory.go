package models

import (
  "time"
)

type Bookcategory struct {
	Bookcategoryid int `gorm:"column:bookcategoryid;primaryKey" json:"bookcategoryid"`
	Categoryname string `gorm:"column:categoryname" json:"categoryname"`
	Categorycode string `gorm:"column:categorycode" json:"categorycode"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookcategory) TableName() string {
	return "bookcategory"
}
