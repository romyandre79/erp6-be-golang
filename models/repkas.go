package models

import (
	"time"
)

type Repkas struct {
	Repkasid    int       `gorm:"column:repkasid;primaryKey" json:"repkasid"`
	Companyid   int       `gorm:"column:companyid" json:"companyid"`
	Nourut      int       `gorm:"column:nourut" json:"nourut"`
	Accountcode string    `gorm:"column:accountcode" json:"accountcode"`
	Plantcode   string    `gorm:"column:plantcode" json:"plantcode"`
	Keterangan  string    `gorm:"column:keterangan" json:"keterangan"`
	Iscount     int8      `gorm:"column:iscount" json:"iscount"`
	Counttype   int8      `gorm:"column:counttype" json:"counttype"`
	Isbold      int8      `gorm:"column:isbold" json:"isbold"`
	Updatedate  time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company     Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Repkas) TableName() string {
	return "repkas"
}
