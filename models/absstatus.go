package models

import (
  "time"
)

type Absstatus struct {
	Absstatusid int `gorm:"column:absstatusid;primaryKey" json:"absstatusid"`
	Shortstat string `gorm:"column:shortstat" json:"shortstat"`
	Longstat string `gorm:"column:longstat" json:"longstat"`
	Isin *int8 `gorm:"column:isin" json:"isin"`
	Priority *int `gorm:"column:priority" json:"priority"`
	Recordstatus *int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Absstatus) TableName() string {
	return "absstatus"
}
