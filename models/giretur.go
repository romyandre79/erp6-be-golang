package models

import (
  "time"
)

type Giretur struct {
	Gireturid int `gorm:"column:gireturid;primaryKey" json:"gireturid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Gireturno *string `gorm:"column:gireturno" json:"gireturno"`
	Fullname *string `gorm:"column:fullname" json:"fullname"`
	Giheaderid int `gorm:"column:giheaderid" json:"giheaderid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Gireturdate time.Time `gorm:"column:gireturdate" json:"gireturdate"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Giheader Giheader `gorm:"foreignKey:giheaderid;references:giheaderid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Giretur) TableName() string {
	return "giretur"
}
