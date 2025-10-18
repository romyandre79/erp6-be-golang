package models

import (
  "time"
)

type Employeetaxdetail struct {
	Employeetaxdetailid int `gorm:"column:employeetaxdetailid;primaryKey" json:"employeetaxdetailid"`
	Employeetaxid int `gorm:"column:employeetaxid" json:"employeetaxid"`
	Wagetypeid int `gorm:"column:wagetypeid" json:"wagetypeid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Reason string `gorm:"column:reason" json:"reason"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Wagetype Wagetype `gorm:"foreignKey:wagetypeid;references:wagetypeid"`
}

func (Employeetaxdetail) TableName() string {
	return "employeetaxdetail"
}
