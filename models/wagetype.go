package models

import (
  "time"
)

type Wagetype struct {
	Wagetypeid int `gorm:"column:wagetypeid;primaryKey" json:"wagetypeid"`
	Wagename string `gorm:"column:wagename" json:"wagename"`
	Ispph int8 `gorm:"column:ispph" json:"ispph"`
	Ispayroll int8 `gorm:"column:ispayroll" json:"ispayroll"`
	Isprint int8 `gorm:"column:isprint" json:"isprint"`
	Percentage float64 `gorm:"column:percentage" json:"percentage"`
	Valuemax float64 `gorm:"column:valuemax" json:"valuemax"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Isrutin int8 `gorm:"column:isrutin" json:"isrutin"`
	Isleft int8 `gorm:"column:isleft" json:"isleft"`
	Paidbycompany int8 `gorm:"column:paidbycompany" json:"paidbycompany"`
	Pphbycompany int8 `gorm:"column:pphbycompany" json:"pphbycompany"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Wagetype) TableName() string {
	return "wagetype"
}
