package models

import (
  "time"
)

type Education struct {
	Educationid int `gorm:"column:educationid;primaryKey" json:"educationid"`
	Educationname string `gorm:"column:educationname" json:"educationname"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Education) TableName() string {
	return "education"
}
