package models

import (
  "time"
)

type Employeekonseling struct {
	Employeekonselingid int `gorm:"column:employeekonselingid;primaryKey" json:"employeekonselingid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Konselingdate time.Time `gorm:"column:konselingdate" json:"konselingdate"`
	Expireddate time.Time `gorm:"column:expireddate" json:"expireddate"`
	Description string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeekonseling) TableName() string {
	return "employeekonseling"
}
