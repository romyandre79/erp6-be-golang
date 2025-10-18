package models

import (

)

type Kutype struct {
	Kutypeid int `gorm:"column:kutypeid;primaryKey" json:"kutypeid"`
	Kutypename string `gorm:"column:kutypename" json:"kutypename"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
}

func (Kutype) TableName() string {
	return "kutype"
}
