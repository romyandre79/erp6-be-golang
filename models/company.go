package models

import (
	"time"
)

type Company struct {
	Companyid     int       `gorm:"column:companyid;primaryKey" json:"companyid"`
	Companyname   string    `gorm:"column:companyname" json:"companyname"`
	Companycode   string    `gorm:"column:companycode" json:"companycode"`
	Address       string    `gorm:"column:address" json:"address"`
	Cityid        int       `gorm:"column:cityid" json:"cityid"`
	Zipcode       string    `gorm:"column:zipcode" json:"zipcode"`
	Taxno         string    `gorm:"column:taxno" json:"taxno"`
	Currencyid    int       `gorm:"column:currencyid" json:"currencyid"`
	Faxno         string    `gorm:"column:faxno" json:"faxno"`
	Phoneno       string    `gorm:"column:phoneno" json:"phoneno"`
	Webaddress    string    `gorm:"column:webaddress" json:"webaddress"`
	Email         string    `gorm:"column:email" json:"email"`
	Leftlogofile  string    `gorm:"column:leftlogofile" json:"leftlogofile"`
	Rightlogofile string    `gorm:"column:rightlogofile" json:"rightlogofile"`
	Isholding     int8      `gorm:"column:isholding" json:"isholding"`
	Billto        string    `gorm:"column:billto" json:"billto"`
	Recordstatus  int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate    time.Time `gorm:"column:updatedate" json:"updatedate"`
	City          City      `gorm:"foreignKey:cityid;references:cityid"`
	Currency      Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
}

func (Company) TableName() string {
	return "company"
}
