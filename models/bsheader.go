package models

import (
	"time"
)

type Bsheader struct {
	Bsheaderid   int       `gorm:"column:bsheaderid;primaryKey" json:"bsheaderid"`
	Plantid      int       `gorm:"column:plantid" json:"plantid"`
	Slocid       int       `gorm:"column:slocid" json:"slocid"`
	Bsdate       time.Time `gorm:"column:bsdate" json:"bsdate"`
	Bsheaderno   string    `gorm:"column:bsheaderno" json:"bsheaderno"`
	Headernote   string    `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname   string    `gorm:"column:statusname" json:"statusname"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant        Plant     `gorm:"foreignKey:plantid;references:plantid"`
}

func (Bsheader) TableName() string {
	return "bsheader"
}
