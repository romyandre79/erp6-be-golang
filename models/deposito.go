package models

import (
  "time"
)

type Deposito struct {
	Depositoid int `gorm:"column:depositoid;primaryKey" json:"depositoid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Depositono *string `gorm:"column:depositono" json:"depositono"`
	Depositodate time.Time `gorm:"column:depositodate" json:"depositodate"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Jenisdepositoid int `gorm:"column:jenisdepositoid" json:"jenisdepositoid"`
	Jumlah float64 `gorm:"column:jumlah" json:"jumlah"`
	Tenor int `gorm:"column:tenor" json:"tenor"`
	Bunga int `gorm:"column:bunga" json:"bunga"`
	Acckas int `gorm:"column:acckas" json:"acckas"`
	Accdeposito int `gorm:"column:accdeposito" json:"accdeposito"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Deposito) TableName() string {
	return "deposito"
}
