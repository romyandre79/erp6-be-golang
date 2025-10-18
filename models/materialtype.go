package models

import (
  "time"
)

type Materialtype struct {
	Materialtypeid int `gorm:"column:materialtypeid;primaryKey" json:"materialtypeid"`
	Materialtypecode string `gorm:"column:materialtypecode" json:"materialtypecode"`
	Description string `gorm:"column:description" json:"description"`
	Fg int8 `gorm:"column:fg" json:"fg"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Materialtype) TableName() string {
	return "materialtype"
}
