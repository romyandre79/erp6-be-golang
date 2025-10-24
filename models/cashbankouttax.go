package models

import (
	"time"
)

type Cashbankouttax struct {
	Cashbankouttaxid int       `gorm:"column:cashbankouttaxid;primaryKey" json:"cashbankouttaxid"`
	Cashbankoutid    int       `gorm:"column:cashbankoutid" json:"cashbankoutid"`
	Taxid            int       `gorm:"column:taxid" json:"taxid"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax              Tax       `gorm:"foreignKey:taxid;references:taxid"`
}

func (Cashbankouttax) TableName() string {
	return "cashbankouttax"
}
