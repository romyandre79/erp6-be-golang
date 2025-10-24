package models

import (
  "time"
)

type Transstructuredet struct {
	Transstructuredetid int `gorm:"column:transstructuredetid;primaryKey" json:"transstructuredetid"`
	Transstructureid int `gorm:"column:transstructureid" json:"transstructureid"`
	Employeeid int `gorm:"column:employeeid" json:"employeeid"`
	Orgstructureid int `gorm:"column:orgstructureid" json:"orgstructureid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Transstructuredet) TableName() string {
	return "transstructuredet"
}
