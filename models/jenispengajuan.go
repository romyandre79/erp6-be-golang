package models

import (
  "time"
)

type Jenispengajuan struct {
	Jenispengajuanid int `gorm:"column:jenispengajuanid;primaryKey" json:"jenispengajuanid"`
	Kodejenispengajuan string `gorm:"column:kodejenispengajuan" json:"kodejenispengajuan"`
	Namajenispengajuan string `gorm:"column:namajenispengajuan" json:"namajenispengajuan"`
	Isfix int8 `gorm:"column:isfix" json:"isfix"`
	Lamaangsuran int8 `gorm:"column:lamaangsuran" json:"lamaangsuran"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Jenispengajuan) TableName() string {
	return "jenispengajuan"
}
