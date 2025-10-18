package models

import (
  "time"
)

type Snrodet struct {
	Snrodid int `gorm:"column:snrodid;primaryKey" json:"snrodid"`
	Plantid *int `gorm:"column:plantid" json:"plantid"`
	Snroid int `gorm:"column:snroid" json:"snroid"`
	Curdd *int `gorm:"column:curdd" json:"curdd"`
	Curmm *int `gorm:"column:curmm" json:"curmm"`
	Curyy *int `gorm:"column:curyy" json:"curyy"`
	Curvalue int `gorm:"column:curvalue" json:"curvalue"`
	Companyid *int `gorm:"column:companyid" json:"companyid"`
	Accountid *int `gorm:"column:accountid" json:"accountid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Snrodet) TableName() string {
	return "snrodet"
}
