package models

import (
	"time"
)

type Absrule struct {
	Absruleid     int         `gorm:"column:absruleid;primaryKey" json:"absruleid"`
	Absscheduleid int         `gorm:"column:absscheduleid" json:"absscheduleid"`
	Difftimein    string      `gorm:"column:difftimein" json:"difftimein"`
	Difftimeout   string      `gorm:"column:difftimeout" json:"difftimeout"`
	Absstatusid   int         `gorm:"column:absstatusid" json:"absstatusid"`
	Recordstatus  int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Absschedule   Absschedule `gorm:"foreignKey:absscheduleid;references:absscheduleid"`
	Absstatus     Absstatus   `gorm:"foreignKey:absstatusid;references:absstatusid"`
}

func (Absrule) TableName() string {
	return "absrule"
}
