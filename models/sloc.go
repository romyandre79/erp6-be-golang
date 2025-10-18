package models

import (
  "time"
)

type Sloc struct {
	Slocid int `gorm:"column:slocid;primaryKey" json:"slocid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Sloccode string `gorm:"column:sloccode" json:"sloccode"`
	Description string `gorm:"column:description" json:"description"`
	Isprd int8 `gorm:"column:isprd" json:"isprd"`
	Isbb int8 `gorm:"column:isbb" json:"isbb"`
	Isbj int8 `gorm:"column:isbj" json:"isbj"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Sloc) TableName() string {
	return "sloc"
}
