package models

import (
	"time"
)

type Addresscontact struct {
	Addresscontactid   int         `gorm:"column:addresscontactid;primaryKey" json:"addresscontactid"`
	Contacttypeid      int         `gorm:"column:contacttypeid" json:"contacttypeid"`
	Addressbookid      int         `gorm:"column:addressbookid" json:"addressbookid"`
	Addresscontactname string      `gorm:"column:addresscontactname" json:"addresscontactname"`
	Phoneno            string      `gorm:"column:phoneno" json:"phoneno"`
	Mobilephone        string      `gorm:"column:mobilephone" json:"mobilephone"`
	Ktp                string      `gorm:"column:ktp" json:"ktp"`
	Emailaddress       string      `gorm:"column:emailaddress" json:"emailaddress"`
	Updatedate         time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Addressbook        Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Contacttype        Contacttype `gorm:"foreignKey:contacttypeid;references:contacttypeid"`
}

func (Addresscontact) TableName() string {
	return "addresscontact"
}
