package models

import (
	"time"
)

type Transstock struct {
	Transstockid     int           `gorm:"column:transstockid;primaryKey" json:"transstockid"`
	Transstocktypeid int8          `gorm:"column:transstocktypeid" json:"transstocktypeid"`
	Plantid          int           `gorm:"column:plantid" json:"plantid"`
	Transstockno     string        `gorm:"column:transstockno" json:"transstockno"`
	Slocfromid       int           `gorm:"column:slocfromid" json:"slocfromid"`
	Slocfromcode     string        `gorm:"column:slocfromcode" json:"slocfromcode"`
	Sloctoid         int           `gorm:"column:sloctoid" json:"sloctoid"`
	Sloctocode       string        `gorm:"column:sloctocode" json:"sloctocode"`
	Transstockdate   time.Time     `gorm:"column:transstockdate" json:"transstockdate"`
	Headernote       string        `gorm:"column:headernote" json:"headernote"`
	Formrequestid    int           `gorm:"column:formrequestid" json:"formrequestid"`
	Formrequestno    string        `gorm:"column:formrequestno" json:"formrequestno"`
	Productoutputid  int           `gorm:"column:productoutputid" json:"productoutputid"`
	Productoutputno  string        `gorm:"column:productoutputno" json:"productoutputno"`
	Addressbookid    int           `gorm:"column:addressbookid" json:"addressbookid"`
	Recordstatus     int8          `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname       string        `gorm:"column:statusname" json:"statusname"`
	Isretur          int8          `gorm:"column:isretur" json:"isretur"`
	Updatedate       time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Plant            Plant         `gorm:"foreignKey:plantid;references:plantid"`
	Productoutput    Productoutput `gorm:"foreignKey:productoutputid;references:productoutputid"`
	Slocfrom         Sloc          `gorm:"foreignKey:slocfromid;references:slocid"`
	Slocto           Sloc          `gorm:"foreignKey:sloctoid;references:slocid"`
}

func (Transstock) TableName() string {
	return "transstock"
}
