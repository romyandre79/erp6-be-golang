package models

import (
	"time"
)

type Emergencydet struct {
	Emergencydetid  int           `gorm:"column:emergencydetid;primaryKey" json:"emergencydetid"`
	Emergencyid     int           `gorm:"column:emergencyid" json:"emergencyid"`
	Emergencytypeid int           `gorm:"column:emergencytypeid" json:"emergencytypeid"`
	Employeeid      int           `gorm:"column:employeeid" json:"employeeid"`
	Evaluasi        string        `gorm:"column:evaluasi" json:"evaluasi"`
	Perbaikan       string        `gorm:"column:perbaikan" json:"perbaikan"`
	Penjelasan      string        `gorm:"column:penjelasan" json:"penjelasan"`
	Description     string        `gorm:"column:description" json:"description"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Emergencytype   Emergencytype `gorm:"foreignKey:emergencytypeid;references:emergencytypeid"`
	Employee        Employee      `gorm:"foreignKey:employeeid;references:employeeid"`
}

func (Emergencydet) TableName() string {
	return "emergencydet"
}
