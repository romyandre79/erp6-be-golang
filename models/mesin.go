package models

import (
	"time"
)

type Mesin struct {
	Mesinid          int       `gorm:"column:mesinid;primaryKey" json:"mesinid"`
	Slocid           int       `gorm:"column:slocid" json:"slocid"`
	Kodemesin        string    `gorm:"column:kodemesin" json:"kodemesin"`
	Namamesin        string    `gorm:"column:namamesin" json:"namamesin"`
	Tahunoperasional string    `gorm:"column:tahunoperasional" json:"tahunoperasional"`
	Buatan           string    `gorm:"column:buatan" json:"buatan"`
	Kwh              float64   `gorm:"column:kwh" json:"kwh"`
	Orgpershift      int       `gorm:"column:orgpershift" json:"orgpershift"`
	Shiftperhari     int       `gorm:"column:shiftperhari" json:"shiftperhari"`
	Speedpermin      int       `gorm:"column:speedpermin" json:"speedpermin"`
	Speedperjam      int       `gorm:"column:speedperjam" json:"speedperjam"`
	Rpm              int       `gorm:"column:rpm" json:"rpm"`
	Lebarbahan       float64   `gorm:"column:lebarbahan" json:"lebarbahan"`
	Rpm2             int       `gorm:"column:rpm2" json:"rpm2"`
	Description      string    `gorm:"column:description" json:"description"`
	Recordstatus     int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Sloc             Sloc      `gorm:"foreignKey:slocid;references:slocid"`
}

func (Mesin) TableName() string {
	return "mesin"
}
