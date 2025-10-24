package models

import (
	"time"
)

type Accounttype struct {
	Accounttypeid       int          `gorm:"column:accounttypeid;primaryKey" json:"accounttypeid"`
	Accounttypename     string       `gorm:"column:accounttypename" json:"accounttypename"`
	Parentaccounttypeid int          `gorm:"column:parentaccounttypeid" json:"parentaccounttypeid"`
	Recordstatus        int8         `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate          time.Time    `gorm:"column:updatedate" json:"updatedate"`
	Parentaccounttype   *Accounttype `gorm:"foreignKey:parentaccounttypeid;references:accounttypeid"`
}

func (Accounttype) TableName() string {
	return "accounttype"
}
