package models

import (
  "time"
)

type Taxrfq struct {
	Taxrfqid int `gorm:"column:taxrfqid;primaryKey" json:"taxrfqid"`
	Rfqid int `gorm:"column:rfqid" json:"rfqid"`
	Taxid int `gorm:"column:taxid" json:"taxid"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax Tax `gorm:"foreignKey:taxid;references:taxid"`
}

func (Taxrfq) TableName() string {
	return "taxrfq"
}
