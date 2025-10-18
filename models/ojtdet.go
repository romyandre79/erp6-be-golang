package models

import (
  "time"
)

type Ojtdet struct {
	Ojtdetid int `gorm:"column:ojtdetid;primaryKey" json:"ojtdetid"`
	Ojtid int `gorm:"column:ojtid" json:"ojtid"`
	Criteriaojtid int `gorm:"column:criteriaojtid" json:"criteriaojtid"`
	Ojtval float64 `gorm:"column:ojtval" json:"ojtval"`
	Correctionval float64 `gorm:"column:correctionval" json:"correctionval"`
	Devisiondesc *string `gorm:"column:devisiondesc" json:"devisiondesc"`
	Departdesc *string `gorm:"column:departdesc" json:"departdesc"`
	Totalval float64 `gorm:"column:totalval" json:"totalval"`
	Averageval float64 `gorm:"column:averageval" json:"averageval"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Ojtdet) TableName() string {
	return "ojtdet"
}
