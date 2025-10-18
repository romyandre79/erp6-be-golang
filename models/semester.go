package models

import (

)

type Semester struct {
	Semesterid int `gorm:"column:semesterid;primaryKey" json:"semesterid"`
	Tahunajaranid int `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Semesternum int `gorm:"column:semesternum" json:"semesternum"`
	Semestertext string `gorm:"column:semestertext" json:"semestertext"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Tahunajaran Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Semester) TableName() string {
	return "semester"
}
