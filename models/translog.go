package models

import (
  "time"
)

type Translog struct {
	Translogid int `gorm:"column:translogid;primaryKey" json:"translogid"`
	Username string `gorm:"column:username" json:"username"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Useraction string `gorm:"column:useraction" json:"useraction"`
	Newdata *string `gorm:"column:newdata" json:"newdata"`
	Menuname *string `gorm:"column:menuname" json:"menuname"`
	Tableid *int `gorm:"column:tableid" json:"tableid"`
	Ippublic *string `gorm:"column:ippublic" json:"ippublic"`
	Iplocal *string `gorm:"column:iplocal" json:"iplocal"`
	Lat *string `gorm:"column:lat" json:"lat"`
	Lng *string `gorm:"column:lng" json:"lng"`
}

func (Translog) TableName() string {
	return "translog"
}
