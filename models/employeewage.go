package models

import (
  "time"
)

type Employeewage struct {
	Employeewageid int `gorm:"column:employeewageid;primaryKey" json:"employeewageid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Wagestartperiod time.Time `gorm:"column:wagestartperiod" json:"wagestartperiod"`
	Workhour int `gorm:"column:workhour" json:"workhour"`
	Wagevalue float64 `gorm:"column:wagevalue" json:"wagevalue"`
	Wageendperiod time.Time `gorm:"column:wageendperiod" json:"wageendperiod"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeewage) TableName() string {
	return "employeewage"
}
