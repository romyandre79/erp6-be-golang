package models

import (
  "time"
)

type Booktype struct {
	Booktypeid int `gorm:"column:booktypeid;primaryKey" json:"booktypeid"`
	Booktypecode string `gorm:"column:booktypecode" json:"booktypecode"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Booktype) TableName() string {
	return "booktype"
}
