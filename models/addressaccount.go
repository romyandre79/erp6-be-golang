package models

import (
	"time"
)

type Addressaccount struct {
	Addressaccountid int         `gorm:"column:addressaccountid;primaryKey" json:"addressaccountid"`
	Companyid        int         `gorm:"column:companyid" json:"companyid"`
	Addressbookid    int         `gorm:"column:addressbookid" json:"addressbookid"`
	Accpiutangid     *int        `gorm:"column:accpiutangid" json:"accpiutangid"`
	Acchutangid      *int        `gorm:"column:acchutangid" json:"acchutangid"`
	Updatedate       time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Acchutang        Account     `gorm:"foreignKey:acchutangid;references:accountid"`
	Accpiutang       Account     `gorm:"foreignKey:accpiutangid;references:accountid"`
	Addressbook      Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Company          Company     `gorm:"foreignKey:companyid;references:companyid"`
}

func (Addressaccount) TableName() string {
	return "addressaccount"
}
