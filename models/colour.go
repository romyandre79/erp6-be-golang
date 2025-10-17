package models

import (
  "time"
)

type Colour struct {
	Colourid int `gorm:"column:colourid;primaryKey" json:"colourid"`
	Colourname string `gorm:"column:colourname" json:"colourname"`
	Webcolour string `gorm:"column:webcolour" json:"webcolour"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Colour) TableName() string {
	return "colour"
}
