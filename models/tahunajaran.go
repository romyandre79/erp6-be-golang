package models

import (

)

type Tahunajaran struct {
	Tahunajaranid int `gorm:"column:tahunajaranid;primaryKey" json:"tahunajaranid"`
	Tahunajaranname string `gorm:"column:tahunajaranname" json:"tahunajaranname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
}

func (Tahunajaran) TableName() string {
	return "tahunajaran"
}
