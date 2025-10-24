package models

import (
  "time"
)

type Identitytype struct {
	Identitytypeid int `gorm:"column:identitytypeid;primaryKey" json:"identitytypeid"`
	Identitytypecode string `gorm:"column:identitytypecode" json:"identitytypecode"`
	Identitytypename string `gorm:"column:identitytypename" json:"identitytypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Identitytype) TableName() string {
	return "identitytype"
}
