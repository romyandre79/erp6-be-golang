package models

import (
  "time"
)

type Grade struct {
	Gradeid int `gorm:"column:gradeid;primaryKey" json:"gradeid"`
	Gradename string `gorm:"column:gradename" json:"gradename"`
	Tunjangan float64 `gorm:"column:tunjangan" json:"tunjangan"`
	Salaryperjam float64 `gorm:"column:salaryperjam" json:"salaryperjam"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Grade) TableName() string {
	return "grade"
}
