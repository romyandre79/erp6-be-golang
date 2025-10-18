package models

import (
  "time"
)

type Plant struct {
	Plantid int `gorm:"column:plantid;primaryKey" json:"plantid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Description string `gorm:"column:description" json:"description"`
	Plantaddress string `gorm:"column:plantaddress" json:"plantaddress"`
	Linkmaps *string `gorm:"column:linkmaps" json:"linkmaps"`
	Limitpo float64 `gorm:"column:limitpo" json:"limitpo"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Plant) TableName() string {
	return "plant"
}
