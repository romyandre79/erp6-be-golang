package models

import (
  "time"
)

type Mrp struct {
	Mrpid int `gorm:"column:mrpid;primaryKey" json:"mrpid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Minstock float64 `gorm:"column:minstock" json:"minstock"`
	Reordervalue float64 `gorm:"column:reordervalue" json:"reordervalue"`
	Maxvalue float64 `gorm:"column:maxvalue" json:"maxvalue"`
	Leadtime int `gorm:"column:leadtime" json:"leadtime"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
}

func (Mrp) TableName() string {
	return "mrp"
}
