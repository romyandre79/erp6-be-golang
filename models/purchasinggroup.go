package models

import (
	"time"
)

type Purchasinggroup struct {
	Purchasinggroupid   int       `gorm:"column:purchasinggroupid;primaryKey" json:"purchasinggroupid"`
	Purchasingorgid     int       `gorm:"column:purchasingorgid" json:"purchasingorgid"`
	Purchasinggroupcode string    `gorm:"column:purchasinggroupcode" json:"purchasinggroupcode"`
	Description         string    `gorm:"column:description" json:"description"`
	Recordstatus        int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate          time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Purchasinggroup) TableName() string {
	return "purchasinggroup"
}
