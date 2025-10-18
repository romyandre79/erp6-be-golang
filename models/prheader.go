package models

import (
  "time"
)

type Prheader struct {
	Prheaderid int `gorm:"column:prheaderid;primaryKey" json:"prheaderid"`
	Prno *string `gorm:"column:prno" json:"prno"`
	Prdate time.Time `gorm:"column:prdate" json:"prdate"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Formrequestid int `gorm:"column:formrequestid" json:"formrequestid"`
	Slocfromid int `gorm:"column:slocfromid" json:"slocfromid"`
	Isjasa int8 `gorm:"column:isjasa" json:"isjasa"`
	Requestedbyid int `gorm:"column:requestedbyid" json:"requestedbyid"`
	Description *string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Requestedby Requestedby `gorm:"foreignKey:requestedbyid;references:requestedbyid"`
	Sloc Sloc `gorm:"foreignKey:slocfromid;references:slocid"`
}

func (Prheader) TableName() string {
	return "prheader"
}
