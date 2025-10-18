package models

import (
  "time"
)

type Jobs struct {
	Jobsid int `gorm:"column:jobsid;primaryKey" json:"jobsid"`
	Orgstructureid int `gorm:"column:orgstructureid" json:"orgstructureid"`
	Jobdesc string `gorm:"column:jobdesc" json:"jobdesc"`
	Qualification string `gorm:"column:qualification" json:"qualification"`
	Positionid int `gorm:"column:positionid" json:"positionid"`
	Isopen int8 `gorm:"column:isopen" json:"isopen"`
	Vacancycount int `gorm:"column:vacancycount" json:"vacancycount"`
	Startdate *time.Time `gorm:"column:startdate" json:"startdate"`
	Enddate *time.Time `gorm:"column:enddate" json:"enddate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Orgstructure Orgstructure `gorm:"foreignKey:orgstructureid;references:orgstructureid"`
	Position Position `gorm:"foreignKey:positionid;references:positionid"`
}

func (Jobs) TableName() string {
	return "jobs"
}
