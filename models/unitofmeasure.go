package models

import (
  "time"
)

type Unitofmeasure struct {
	Unitofmeasureid int `gorm:"column:unitofmeasureid;primaryKey" json:"unitofmeasureid"`
	Uomcode string `gorm:"column:uomcode" json:"uomcode"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Unitofmeasure) TableName() string {
	return "unitofmeasure"
}
