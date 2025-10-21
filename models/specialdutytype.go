package models

import (
	"time"
)

type Specialdutytype struct {
	Specialdutytypeid   int       `gorm:"column:specialdutytypeid;primaryKey" json:"specialdutytypeid"`
	Specialdutytypename string    `gorm:"column:specialdutytypename" json:"specialdutytypename"`
	Description         string    `gorm:"column:description" json:"description"`
	Recordstatus        int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate          time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Specialdutytype) TableName() string {
	return "specialdutytype"
}
