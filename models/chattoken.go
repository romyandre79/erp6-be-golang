package models

import (
  "time"
)

type Chattoken struct {
	Chattokenid int `gorm:"column:chattokenid;primaryKey" json:"chattokenid"`
	Token string `gorm:"column:token" json:"token"`
	Tokentype int8 `gorm:"column:tokentype" json:"tokentype"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Chattoken) TableName() string {
	return "chattoken"
}
