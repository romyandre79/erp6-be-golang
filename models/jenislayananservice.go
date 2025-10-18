package models

import (
  "time"
)

type Jenislayananservice struct {
	Jenislayananserviceid int `gorm:"column:jenislayananserviceid;primaryKey" json:"jenislayananserviceid"`
	Jenislayananservicename string `gorm:"column:jenislayananservicename" json:"jenislayananservicename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Jenislayananservice) TableName() string {
	return "jenislayananservice"
}
