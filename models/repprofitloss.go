package models

import (
  "time"
)

type Repprofitloss struct {
	Repprofitlossid int `gorm:"column:repprofitlossid;primaryKey" json:"repprofitlossid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Nourut *int8 `gorm:"column:nourut" json:"nourut"`
	Accountcode *string `gorm:"column:accountcode" json:"accountcode"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Iscount int8 `gorm:"column:iscount" json:"iscount"`
	Counttype int8 `gorm:"column:counttype" json:"counttype"`
	Sourcemenu *string `gorm:"column:sourcemenu" json:"sourcemenu"`
	Isbold *int8 `gorm:"column:isbold" json:"isbold"`
	Keterangan string `gorm:"column:keterangan" json:"keterangan"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Repprofitloss) TableName() string {
	return "repprofitloss"
}
