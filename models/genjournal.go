package models

import (
	"time"
)

type Genjournal struct {
	Genjournalid int       `gorm:"column:genjournalid;primaryKey" json:"genjournalid"`
	Companyid    int       `gorm:"column:companyid" json:"companyid"`
	Companycode  string    `gorm:"column:companycode" json:"companycode"`
	Plantid      int       `gorm:"column:plantid" json:"plantid"`
	Plantcode    string    `gorm:"column:plantcode" json:"plantcode"`
	Journalno    string    `gorm:"column:journalno" json:"journalno"`
	Referenceno  string    `gorm:"column:referenceno" json:"referenceno"`
	Journaldate  time.Time `gorm:"column:journaldate" json:"journaldate"`
	Postdate     time.Time `gorm:"column:postdate" json:"postdate"`
	Journalnote  string    `gorm:"column:journalnote" json:"journalnote"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname   string    `gorm:"column:statusname" json:"statusname"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company      Company   `gorm:"foreignKey:companyid;references:companyid"`
	Plant        Plant     `gorm:"foreignKey:plantid;references:plantid"`
}

func (Genjournal) TableName() string {
	return "genjournal"
}
