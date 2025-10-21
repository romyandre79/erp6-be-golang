package models

import (
	"time"
)

type Payrollperiod struct {
	Payrollperiodid   int            `gorm:"column:payrollperiodid;primaryKey" json:"payrollperiodid"`
	Payrollperiodname string         `gorm:"column:payrollperiodname" json:"payrollperiodname"`
	Startdate         time.Time      `gorm:"column:startdate" json:"startdate"`
	Enddate           time.Time      `gorm:"column:enddate" json:"enddate"`
	Parentperiodid    int            `gorm:"column:parentperiodid" json:"parentperiodid"`
	Recordstatus      int8           `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate        time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Parentperiod      *Payrollperiod `gorm:"foreignKey:parentperiodid;references:payrollperiodid"`
}

func (Payrollperiod) TableName() string {
	return "payrollperiod"
}
