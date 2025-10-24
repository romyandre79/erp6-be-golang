package models

import (
	"time"
)

type Nilaitugas struct {
	Nilaitugasid   int       `gorm:"column:nilaitugasid;primaryKey" json:"nilaitugasid"`
	Tugasid        int       `gorm:"column:tugasid" json:"tugasid"`
	Addressbookid  int       `gorm:"column:addressbookid" json:"addressbookid"`
	File           string    `gorm:"column:file" json:"file"`
	Content        string    `gorm:"column:content" json:"content"`
	Nilai          float64   `gorm:"column:nilai" json:"nilai"`
	Nilaitugasdate time.Time `gorm:"column:nilaitugasdate" json:"nilaitugasdate"`
}

func (Nilaitugas) TableName() string {
	return "nilaitugas"
}
