package models

import (

)

type Tugaskelas struct {
	Tugaskelasid int `gorm:"column:tugaskelasid;primaryKey" json:"tugaskelasid"`
	Tugasid int `gorm:"column:tugasid" json:"tugasid"`
	Kelasid int `gorm:"column:kelasid" json:"kelasid"`
	Tugas Tugas `gorm:"foreignKey:tugasid;references:tugasid"`
}

func (Tugaskelas) TableName() string {
	return "tugaskelas"
}
