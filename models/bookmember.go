package models

import (
  "time"
)

type Bookmember struct {
	Bookmemberid int `gorm:"column:bookmemberid;primaryKey" json:"bookmemberid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Bookmemberno *string `gorm:"column:bookmemberno" json:"bookmemberno"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Membertypeid int `gorm:"column:membertypeid" json:"membertypeid"`
	Registerdate time.Time `gorm:"column:registerdate" json:"registerdate"`
	Expiredate time.Time `gorm:"column:expiredate" json:"expiredate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
	Bookmembertype Bookmembertype `gorm:"foreignKey:membertypeid;references:bookmembertypeid"`
}

func (Bookmember) TableName() string {
	return "bookmember"
}
