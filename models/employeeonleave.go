package models

import (
  "time"
)

type Employeeonleave struct {
	Employeeonleaveid int `gorm:"column:employeeonleaveid;primaryKey" json:"employeeonleaveid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Onleaveid int `gorm:"column:onleaveid" json:"onleaveid"`
	Periodefrom time.Time `gorm:"column:periodefrom" json:"periodefrom"`
	Periodeto time.Time `gorm:"column:periodeto" json:"periodeto"`
	Onleavequota int `gorm:"column:onleavequota" json:"onleavequota"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Onleavetype Onleavetype `gorm:"foreignKey:onleaveid;references:onleavetypeid"`
}

func (Employeeonleave) TableName() string {
	return "employeeonleave"
}
