package models

import (
  "time"
)

type Budget struct {
	Budgetid int `gorm:"column:budgetid;primaryKey" json:"budgetid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Budgetamount float64 `gorm:"column:budgetamount" json:"budgetamount"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Amountused *float64 `gorm:"column:amountused" json:"amountused"`
	Amountreqused *float64 `gorm:"column:amountreqused" json:"amountreqused"`
	Isstrict *int8 `gorm:"column:isstrict" json:"isstrict"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Sloc Sloc `gorm:"foreignKey:slocid;references:slocid"`
}

func (Budget) TableName() string {
	return "budget"
}
