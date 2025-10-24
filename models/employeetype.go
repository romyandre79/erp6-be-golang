package models

import (
	"time"
)

type Employeetype struct {
	Employeetypeid   int       `gorm:"column:employeetypeid;primaryKey" json:"employeetypeid"`
	Employeetypename string    `gorm:"column:employeetypename" json:"employeetypename"`
	Snroid           int       `gorm:"column:snroid" json:"snroid"`
	Sicksnroid       int       `gorm:"column:sicksnroid" json:"sicksnroid"`
	Sickstatusid     int       `gorm:"column:sickstatusid" json:"sickstatusid"`
	Recordstatus     int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Sicksnro         Snro      `gorm:"foreignKey:sicksnroid;references:snroid"`
	Snro             Snro      `gorm:"foreignKey:snroid;references:snroid"`
}

func (Employeetype) TableName() string {
	return "employeetype"
}
