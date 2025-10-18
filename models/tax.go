package models

import (
  "time"
)

type Tax struct {
	Taxid int `gorm:"column:taxid;primaryKey" json:"taxid"`
	Taxcode string `gorm:"column:taxcode" json:"taxcode"`
	Taxvalue float64 `gorm:"column:taxvalue" json:"taxvalue"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Tax) TableName() string {
	return "tax"
}
