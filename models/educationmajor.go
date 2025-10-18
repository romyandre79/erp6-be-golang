package models

import (
  "time"
)

type Educationmajor struct {
	Educationmajorid int `gorm:"column:educationmajorid;primaryKey" json:"educationmajorid"`
	Educationmajorname string `gorm:"column:educationmajorname" json:"educationmajorname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Educationmajor) TableName() string {
	return "educationmajor"
}
