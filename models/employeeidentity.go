package models

import (
  "time"
)

type Employeeidentity struct {
	Employeeidentityid int `gorm:"column:employeeidentityid;primaryKey" json:"employeeidentityid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Identitytypeid int `gorm:"column:identitytypeid" json:"identitytypeid"`
	Identityname string `gorm:"column:identityname" json:"identityname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Identitytype Identitytype `gorm:"foreignKey:identitytypeid;references:identitytypeid"`
}

func (Employeeidentity) TableName() string {
	return "employeeidentity"
}
