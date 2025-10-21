package models

import (
	"time"
)

type Workaccident struct {
	Workaccidentid     int              `gorm:"column:workaccidentid;primaryKey" json:"workaccidentid"`
	Workaccidentdate   time.Time        `gorm:"column:workaccidentdate" json:"workaccidentdate"`
	Workaccidentno     string           `gorm:"column:workaccidentno" json:"workaccidentno"`
	Plantid            int              `gorm:"column:plantid" json:"plantid"`
	Workaccidenttypeid int              `gorm:"column:workaccidenttypeid" json:"workaccidenttypeid"`
	Employeeid         int              `gorm:"column:employeeid" json:"employeeid"`
	Orgstructureid     int              `gorm:"column:orgstructureid" json:"orgstructureid"`
	Accidentreport     string           `gorm:"column:accidentreport" json:"accidentreport"`
	Description        string           `gorm:"column:description" json:"description"`
	Recordstatus       int8             `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname         string           `gorm:"column:statusname" json:"statusname"`
	Updatedate         time.Time        `gorm:"column:updatedate" json:"updatedate"`
	Employee           Employee         `gorm:"foreignKey:employeeid;references:employeeid"`
	Orgstructure       Orgstructure     `gorm:"foreignKey:orgstructureid;references:orgstructureid"`
	Plant              Plant            `gorm:"foreignKey:plantid;references:plantid"`
	Workaccidenttype   Workaccidenttype `gorm:"foreignKey:workaccidenttypeid;references:workaccidenttypeid"`
}

func (Workaccident) TableName() string {
	return "workaccident"
}
