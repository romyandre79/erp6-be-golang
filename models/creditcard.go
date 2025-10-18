package models

import (
  "time"
)

type Creditcard struct {
	Creditcardid int `gorm:"column:creditcardid;primaryKey" json:"creditcardid"`
	Creditcardname string `gorm:"column:creditcardname" json:"creditcardname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Creditcard) TableName() string {
	return "creditcard"
}
