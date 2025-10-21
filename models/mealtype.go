package models

import (
	"time"
)

type Mealtype struct {
	Mealtypeid   int       `gorm:"column:mealtypeid;primaryKey" json:"mealtypeid"`
	Mealtypename string    `gorm:"column:mealtypename" json:"mealtypename"`
	Description  string    `gorm:"column:description" json:"description"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Mealtype) TableName() string {
	return "mealtype"
}
