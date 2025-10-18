package models

import (
  "time"
)

type Familyrelation struct {
	Familyrelationid int `gorm:"column:familyrelationid;primaryKey" json:"familyrelationid"`
	Familyrelationname string `gorm:"column:familyrelationname" json:"familyrelationname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Familyrelation) TableName() string {
	return "familyrelation"
}
