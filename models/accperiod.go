package models

import (
  "time"
)

type Accperiod struct {
	Accperiodid int `gorm:"column:accperiodid;primaryKey" json:"accperiodid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Perioddate time.Time `gorm:"column:perioddate" json:"perioddate"`
	Accstatus int8 `gorm:"column:accstatus" json:"accstatus"`
	Invstatus int8 `gorm:"column:invstatus" json:"invstatus"`
	Purstatus int8 `gorm:"column:purstatus" json:"purstatus"`
	Ordstatus int8 `gorm:"column:ordstatus" json:"ordstatus"`
	Prdstatus int8 `gorm:"column:prdstatus" json:"prdstatus"`
	Hrstatus int8 `gorm:"column:hrstatus" json:"hrstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Accperiod) TableName() string {
	return "accperiod"
}
