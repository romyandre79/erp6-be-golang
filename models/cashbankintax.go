package models

import (
  "time"
)

type Cashbankintax struct {
	Cashbankintaxid int `gorm:"column:cashbankintaxid;primaryKey" json:"cashbankintaxid"`
	Cashbankinid int `gorm:"column:cashbankinid" json:"cashbankinid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Cashbankintax) TableName() string {
	return "cashbankintax"
}
