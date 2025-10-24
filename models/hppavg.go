package models

import (

)

type Hppavg struct {
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Transdate int8 `gorm:"column:transdate" json:"transdate"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Rp float64 `gorm:"column:rp" json:"rp"`
	Qtyperrp float64 `gorm:"column:qtyperrp" json:"qtyperrp"`
}

func (Hppavg) TableName() string {
	return "hppavg"
}
