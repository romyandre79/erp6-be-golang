package models

import (
  "time"
)

type Position struct {
	Positionid int `gorm:"column:positionid;primaryKey" json:"positionid"`
	Positionname string `gorm:"column:positionname" json:"positionname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Position) TableName() string {
	return "position"
}
