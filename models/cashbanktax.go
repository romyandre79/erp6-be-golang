package models

import (
  "time"
)

type Cashbanktax struct {
	Cashbanktaxid int `gorm:"column:cashbanktaxid;primaryKey" json:"cashbanktaxid"`
	Cashbankid int `gorm:"column:cashbankid" json:"cashbankid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Nobuktipotong string `gorm:"column:nobuktipotong" json:"nobuktipotong"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Cashbanktax) TableName() string {
	return "cashbanktax"
}
