package models

import (
  "time"
)

type Repbiayaumum struct {
	Repbiayaumumid int `gorm:"column:repbiayaumumid;primaryKey" json:"repbiayaumumid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Accountcode string `gorm:"column:accountcode" json:"accountcode"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Nourut string `gorm:"column:nourut" json:"nourut"`
	Isbold int8 `gorm:"column:isbold" json:"isbold"`
	Iscount int8 `gorm:"column:iscount" json:"iscount"`
	Counttypeid int `gorm:"column:counttypeid" json:"counttypeid"`
	Keterangan string `gorm:"column:keterangan" json:"keterangan"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Repbiayaumum) TableName() string {
	return "repbiayaumum"
}
