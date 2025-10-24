package models

import (
	"time"
)

type Rephpp struct {
	Rephppid    int       `gorm:"column:rephppid;primaryKey" json:"rephppid"`
	Companyid   int       `gorm:"column:companyid" json:"companyid"`
	Nourut      int       `gorm:"column:nourut" json:"nourut"`
	Accountcode string    `gorm:"column:accountcode" json:"accountcode"`
	Plantcode   string    `gorm:"column:plantcode" json:"plantcode"`
	Keterangan  string    `gorm:"column:keterangan" json:"keterangan"`
	Iscount     int8      `gorm:"column:iscount" json:"iscount"`
	Counttype   int8      `gorm:"column:counttype" json:"counttype"`
	Isbold      int8      `gorm:"column:isbold" json:"isbold"`
	Updatedate  time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Rephpp) TableName() string {
	return "rephpp"
}
