package models

import (
  "time"
)

type Romawi struct {
	Romawiid int `gorm:"column:romawiid;primaryKey" json:"romawiid"`
	Monthcal int `gorm:"column:monthcal" json:"monthcal"`
	Monthrm string `gorm:"column:monthrm" json:"monthrm"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Romawi) TableName() string {
	return "romawi"
}
