package models

import (
  "time"
)

type Lamaangsuran struct {
	Lamaangsuranid int `gorm:"column:lamaangsuranid;primaryKey" json:"lamaangsuranid"`
	Kodelamaangsuran string `gorm:"column:kodelamaangsuran" json:"kodelamaangsuran"`
	Jangkawaktu int `gorm:"column:jangkawaktu" json:"jangkawaktu"`
	Isauto int8 `gorm:"column:isauto" json:"isauto"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Lamaangsuran) TableName() string {
	return "lamaangsuran"
}
