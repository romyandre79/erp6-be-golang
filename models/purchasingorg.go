package models

import (
  "time"
)

type Purchasingorg struct {
	Purchasingorgid int `gorm:"column:purchasingorgid;primaryKey" json:"purchasingorgid"`
	Purchasingorgcode string `gorm:"column:purchasingorgcode" json:"purchasingorgcode"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Purchasingorg) TableName() string {
	return "purchasingorg"
}
