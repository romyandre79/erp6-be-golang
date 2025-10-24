package models

import (
	"time"
)

type Account struct {
	Accountid       int         `gorm:"column:accountid;primaryKey" json:"accountid"`
	Companyid       int         `gorm:"column:companyid" json:"companyid"`
	Accountcode     string      `gorm:"column:accountcode" json:"accountcode"`
	Accountname     string      `gorm:"column:accountname" json:"accountname"`
	Parentaccountid int         `gorm:"column:parentaccountid" json:"parentaccountid"`
	Currencyid      int         `gorm:"column:currencyid" json:"currencyid"`
	Accounttypeid   int         `gorm:"column:accounttypeid" json:"accounttypeid"`
	Isdebit         int8        `gorm:"column:isdebit" json:"isdebit"`
	Recordstatus    int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate      time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Company         Company     `gorm:"foreignKey:companyid;references:companyid"`
	Currency        Currency    `gorm:"foreignKey:currencyid;references:currencyid"`
	Parentaccount   *Account    `gorm:"foreignKey:parentaccountid;references:accountid"`
	Accounttype     Accounttype `gorm:"foreignKey:accounttypeid;references:accounttypeid"`
}

func (Account) TableName() string {
	return "account"
}
