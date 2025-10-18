package models

import (
  "time"
)

type Cashbank struct {
	Cashbankid int `gorm:"column:cashbankid;primaryKey" json:"cashbankid"`
	Cashbankno *string `gorm:"column:cashbankno" json:"cashbankno"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Cashbankdate time.Time `gorm:"column:cashbankdate" json:"cashbankdate"`
	Isin int8 `gorm:"column:isin" json:"isin"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Accountid int `gorm:"column:accountid" json:"accountid"`
	Receiptno string `gorm:"column:receiptno" json:"receiptno"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Headernote string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Account Account `gorm:"foreignKey:accountid;references:accountid"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Cashbank) TableName() string {
	return "cashbank"
}
