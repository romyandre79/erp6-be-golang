package models

import (
  "time"
)

type Notagrrgrr struct {
	Notagrrgrrid int `gorm:"column:notagrrgrrid;primaryKey" json:"notagrrgrrid"`
	Notagrrid int `gorm:"column:notagrrid" json:"notagrrid"`
	Grreturid int `gorm:"column:grreturid" json:"grreturid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Grretur Grretur `gorm:"foreignKey:grreturid;references:grreturid"`
}

func (Notagrrgrr) TableName() string {
	return "notagrrgrr"
}
