package models

import (
  "time"
)

type Employeetax struct {
	Employeetaxid int `gorm:"column:employeetaxid;primaryKey" json:"employeetaxid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Taxstartperiod time.Time `gorm:"column:taxstartperiod" json:"taxstartperiod"`
	Workhour float64 `gorm:"column:workhour" json:"workhour"`
	Taxvalue float64 `gorm:"column:taxvalue" json:"taxvalue"`
	Taxendperiod time.Time `gorm:"column:taxendperiod" json:"taxendperiod"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Employeetax) TableName() string {
	return "employeetax"
}
