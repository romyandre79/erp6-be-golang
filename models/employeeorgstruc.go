package models

import (
  "time"
)

type Employeeorgstruc struct {
	Employeeorgstrucid int `gorm:"column:employeeorgstrucid;primaryKey" json:"employeeorgstrucid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Orgstructureid int `gorm:"column:orgstructureid" json:"orgstructureid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Orgstructure Orgstructure `gorm:"foreignKey:orgstructureid;references:orgstructureid"`
}

func (Employeeorgstruc) TableName() string {
	return "employeeorgstruc"
}
