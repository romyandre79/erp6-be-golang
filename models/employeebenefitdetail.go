package models

import (
  "time"
)

type Employeebenefitdetail struct {
	Employeebenefitdetailid int `gorm:"column:employeebenefitdetailid;primaryKey" json:"employeebenefitdetailid"`
	Employeebenefitid int `gorm:"column:employeebenefitid" json:"employeebenefitid"`
	Wagetypeid int `gorm:"column:wagetypeid" json:"wagetypeid"`
	Startdate time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Isfinal int8 `gorm:"column:isfinal" json:"isfinal"`
	Reason *string `gorm:"column:reason" json:"reason"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Wagetype Wagetype `gorm:"foreignKey:wagetypeid;references:wagetypeid"`
}

func (Employeebenefitdetail) TableName() string {
	return "employeebenefitdetail"
}
