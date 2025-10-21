package models

import (
	"time"
)

type Grjurnal struct {
	Grjurnalid int       `gorm:"column:grjurnalid;primaryKey" json:"grjurnalid"`
	Grheaderid int       `gorm:"column:grheaderid" json:"grheaderid"`
	Accountid  int       `gorm:"column:accountid" json:"accountid"`
	Debet      float64   `gorm:"column:debet" json:"debet"`
	Credit     float64   `gorm:"column:credit" json:"credit"`
	Currencyid int       `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue  int       `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account    Account   `gorm:"foreignKey:accountid;references:accountid"`
	Currency   Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
	Grheader   Grheader  `gorm:"foreignKey:grheaderid;references:grheaderid"`
}

func (Grjurnal) TableName() string {
	return "grjurnal"
}
