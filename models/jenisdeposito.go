package models

import (
  "time"
)

type Jenisdeposito struct {
	Jenisdepositoid int `gorm:"column:jenisdepositoid;primaryKey" json:"jenisdepositoid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Namadeposito string `gorm:"column:namadeposito" json:"namadeposito"`
	Jumlah float64 `gorm:"column:jumlah" json:"jumlah"`
	Bunga float64 `gorm:"column:bunga" json:"bunga"`
	Fixed int8 `gorm:"column:fixed" json:"fixed"`
	Tenor int `gorm:"column:tenor" json:"tenor"`
	Isauto int8 `gorm:"column:isauto" json:"isauto"`
	Acckas int `gorm:"column:acckas" json:"acckas"`
	Accdeposito int `gorm:"column:accdeposito" json:"accdeposito"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Jenisdeposito) TableName() string {
	return "jenisdeposito"
}
