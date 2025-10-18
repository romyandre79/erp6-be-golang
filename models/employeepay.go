package models

import (
  "time"
)

type Employeepay struct {
	Employeepayid int `gorm:"column:employeepayid;primaryKey" json:"employeepayid"`
	Employeepayno *string `gorm:"column:employeepayno" json:"employeepayno"`
	Employeepaydate time.Time `gorm:"column:employeepaydate" json:"employeepaydate"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Employeepiutangid int `gorm:"column:employeepiutangid" json:"employeepiutangid"`
	Nilai float64 `gorm:"column:nilai" json:"nilai"`
	Sisa float64 `gorm:"column:sisa" json:"sisa"`
	Description *string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employeepiutang Employeepiutang `gorm:"foreignKey:employeepiutangid;references:employeepiutangid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Employeepay) TableName() string {
	return "employeepay"
}
