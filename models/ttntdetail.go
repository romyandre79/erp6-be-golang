package models

import (
  "time"
)

type Ttntdetail struct {
	Ttntdetailid int `gorm:"column:ttntdetailid;primaryKey" json:"ttntdetailid"`
	Ttntid int `gorm:"column:ttntid" json:"ttntid"`
	Invoicearid int `gorm:"column:invoicearid" json:"invoicearid"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Ttntdetail) TableName() string {
	return "ttntdetail"
}
