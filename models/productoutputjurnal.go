package models

import (
	"time"
)

type Productoutputjurnal struct {
	Productoutputjurnalid int       `gorm:"column:productoutputjurnalid;primaryKey" json:"productoutputjurnalid"`
	Productid             int       `gorm:"column:productid" json:"productid"`
	Productoutputfgid     int       `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Productoutputid       int       `gorm:"column:productoutputid" json:"productoutputid"`
	Accountid             int       `gorm:"column:accountid" json:"accountid"`
	Accountname           string    `gorm:"column:accountname" json:"accountname"`
	Debit                 float64   `gorm:"column:debit" json:"debit"`
	Credit                float64   `gorm:"column:credit" json:"credit"`
	Currencyid            int       `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue             float64   `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote            string    `gorm:"column:detailnote" json:"detailnote"`
	Updatedate            time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Productoutputjurnal) TableName() string {
	return "productoutputjurnal"
}
