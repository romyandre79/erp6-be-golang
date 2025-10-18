package models

import (
  "time"
)

type Invoiceaptax struct {
	Invoiceaptaxid int `gorm:"column:invoiceaptaxid;primaryKey" json:"invoiceaptaxid"`
	Invoiceapid int `gorm:"column:invoiceapid" json:"invoiceapid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Invoiceaptax) TableName() string {
	return "invoiceaptax"
}
