package models

import (
  "time"
)

type Userinbox struct {
	Userinboxid int `gorm:"column:userinboxid;primaryKey" json:"userinboxid"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Fromuserid int `gorm:"column:fromuserid" json:"fromuserid"`
	Subject string `gorm:"column:subject" json:"subject"`
	Description string `gorm:"column:description" json:"description"`
	Senddate time.Time `gorm:"column:senddate" json:"senddate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Userinbox) TableName() string {
	return "userinbox"
}
