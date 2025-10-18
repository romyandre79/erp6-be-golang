package models

import (
  "time"
)

type Jenisbunga struct {
	Jenisbungaid int `gorm:"column:jenisbungaid;primaryKey" json:"jenisbungaid"`
	Jenisbunganame string `gorm:"column:jenisbunganame" json:"jenisbunganame"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Jenisbunga) TableName() string {
	return "jenisbunga"
}
