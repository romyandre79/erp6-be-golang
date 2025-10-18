package models

import (
  "time"
)

type Taxwageprogressif struct {
	Taxwageprogressifid int `gorm:"column:taxwageprogressifid;primaryKey" json:"taxwageprogressifid"`
	Description string `gorm:"column:description" json:"description"`
	Minvalue float64 `gorm:"column:minvalue" json:"minvalue"`
	Maxvalues float64 `gorm:"column:maxvalues" json:"maxvalues"`
	Valuepercent float64 `gorm:"column:valuepercent" json:"valuepercent"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Taxwageprogressif) TableName() string {
	return "taxwageprogressif"
}
