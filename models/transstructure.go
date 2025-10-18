package models

import (
  "time"
)

type Transstructure struct {
	Transstructureid int `gorm:"column:transstructureid;primaryKey" json:"transstructureid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Docdate time.Time `gorm:"column:docdate" json:"docdate"`
	Docno *string `gorm:"column:docno" json:"docno"`
	Description *string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Transstructure) TableName() string {
	return "transstructure"
}
