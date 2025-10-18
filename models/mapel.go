package models

import (
  "time"
)

type Mapel struct {
	Mapelid int `gorm:"column:mapelid;primaryKey" json:"mapelid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Mapelname string `gorm:"column:mapelname" json:"mapelname"`
	Mapelpic string `gorm:"column:mapelpic" json:"mapelpic"`
	Description string `gorm:"column:description" json:"description"`
	Mapeldesc string `gorm:"column:mapeldesc" json:"mapeldesc"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Mapel) TableName() string {
	return "mapel"
}
