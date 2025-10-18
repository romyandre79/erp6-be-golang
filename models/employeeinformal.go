package models

import (
  "time"
)

type Employeeinformal struct {
	Employeeinformalid int `gorm:"column:employeeinformalid;primaryKey" json:"employeeinformalid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Informalname string `gorm:"column:informalname" json:"informalname"`
	Organizer string `gorm:"column:organizer" json:"organizer"`
	Period string `gorm:"column:period" json:"period"`
	Isdiploma int8 `gorm:"column:isdiploma" json:"isdiploma"`
	Sponsoredby string `gorm:"column:sponsoredby" json:"sponsoredby"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Employeeinformal) TableName() string {
	return "employeeinformal"
}
