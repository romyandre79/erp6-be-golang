package models

import (
  "time"
)

type Grretur struct {
	Grreturid int `gorm:"column:grreturid;primaryKey" json:"grreturid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Grreturno *string `gorm:"column:grreturno" json:"grreturno"`
	Grreturdate time.Time `gorm:"column:grreturdate" json:"grreturdate"`
	Poheaderid int `gorm:"column:poheaderid" json:"poheaderid"`
	Fullname string `gorm:"column:fullname" json:"fullname"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Poheader Poheader `gorm:"foreignKey:poheaderid;references:poheaderid"`
}

func (Grretur) TableName() string {
	return "grretur"
}
