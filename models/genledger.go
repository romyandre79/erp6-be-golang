package models

import (
	"time"
)

type Genledger struct {
	Genledgerid  int       `gorm:"column:genledgerid;primaryKey" json:"genledgerid"`
	Companyid    int       `gorm:"column:companyid" json:"companyid"`
	Companycode  string    `gorm:"column:companycode" json:"companycode"`
	Plantid      int       `gorm:"column:plantid" json:"plantid"`
	Plantcode    string    `gorm:"column:plantcode" json:"plantcode"`
	Accountid    int       `gorm:"column:accountid" json:"accountid"`
	Accountcode  string    `gorm:"column:accountcode" json:"accountcode"`
	Accountname  string    `gorm:"column:accountname" json:"accountname"`
	Genjournalid int       `gorm:"column:genjournalid" json:"genjournalid"`
	Debit        float64   `gorm:"column:debit" json:"debit"`
	Credit       float64   `gorm:"column:credit" json:"credit"`
	Postdate     time.Time `gorm:"column:postdate" json:"postdate"`
	Journaldate  time.Time `gorm:"column:journaldate" json:"journaldate"`
	Currencyid   int       `gorm:"column:currencyid" json:"currencyid"`
	Currencyname string    `gorm:"column:currencyname" json:"currencyname"`
	Symbol       string    `gorm:"column:symbol" json:"symbol"`
	Ratevalue    float64   `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote   string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company      Company   `gorm:"foreignKey:companyid;references:companyid"`
	Currency     Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Genledger) TableName() string {
	return "genledger"
}
