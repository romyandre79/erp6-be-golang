package models

import (
  "time"
)

type Journaldetail struct {
	Journaldetailid int `gorm:"column:journaldetailid;primaryKey" json:"journaldetailid"`
	Genjournalid int `gorm:"column:genjournalid" json:"genjournalid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Debit *float64 `gorm:"column:debit" json:"debit"`
	Credit *float64 `gorm:"column:credit" json:"credit"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote *string `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Journaldetail) TableName() string {
	return "journaldetail"
}
