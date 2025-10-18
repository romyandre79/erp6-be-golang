package models

import (
  "time"
)

type Repbiayapenjualan struct {
	Repbiayapenjualanid int `gorm:"column:repbiayapenjualanid;primaryKey" json:"repbiayapenjualanid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Accountcode string `gorm:"column:accountcode" json:"accountcode"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Nourut string `gorm:"column:nourut" json:"nourut"`
	Isbold int8 `gorm:"column:isbold" json:"isbold"`
	Counttypeid int `gorm:"column:counttypeid" json:"counttypeid"`
	Iscount int8 `gorm:"column:iscount" json:"iscount"`
	Keterangan string `gorm:"column:keterangan" json:"keterangan"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Repbiayapenjualan) TableName() string {
	return "repbiayapenjualan"
}
