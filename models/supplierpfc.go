package models

import (

)

type Supplierpfc struct {
	Supplierpfcid int `gorm:"column:supplierpfcid;primaryKey" json:"supplierpfcid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Periodmonth int `gorm:"column:periodmonth" json:"periodmonth"`
	Periodyear int `gorm:"column:periodyear" json:"periodyear"`
	Contractcount int `gorm:"column:contractcount" json:"contractcount"`
	Contractamount float64 `gorm:"column:contractamount" json:"contractamount"`
	Performance float64 `gorm:"column:performance" json:"performance"`
	Grade string `gorm:"column:grade" json:"grade"`
	Note string `gorm:"column:note" json:"note"`
}

func (Supplierpfc) TableName() string {
	return "supplierpfc"
}
