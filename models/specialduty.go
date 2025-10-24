package models

import (
	"time"
)

type Specialduty struct {
	Specialdutyid     int             `gorm:"column:specialdutyid;primaryKey" json:"specialdutyid"`
	Plantid           int             `gorm:"column:plantid" json:"plantid"`
	Employeeid        int             `gorm:"column:employeeid" json:"employeeid"`
	Specialdutydate   time.Time       `gorm:"column:specialdutydate" json:"specialdutydate"`
	Specialdutyno     string          `gorm:"column:specialdutyno" json:"specialdutyno"`
	Positionid        int             `gorm:"column:positionid" json:"positionid"`
	Specialdutytypeid int             `gorm:"column:specialdutytypeid" json:"specialdutytypeid"`
	Orgstructureid    int             `gorm:"column:orgstructureid" json:"orgstructureid"`
	Description       string          `gorm:"column:description" json:"description"`
	Recordstatus      int8            `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname        string          `gorm:"column:statusname" json:"statusname"`
	Updatedate        time.Time       `gorm:"column:updatedate" json:"updatedate"`
	Employee          Employee        `gorm:"foreignKey:employeeid;references:employeeid"`
	Orgstructure      Orgstructure    `gorm:"foreignKey:orgstructureid;references:orgstructureid"`
	Plant             Plant           `gorm:"foreignKey:plantid;references:plantid"`
	Position          Position        `gorm:"foreignKey:positionid;references:positionid"`
	Specialdutytype   Specialdutytype `gorm:"foreignKey:specialdutytypeid;references:specialdutytypeid"`
}

func (Specialduty) TableName() string {
	return "specialduty"
}
