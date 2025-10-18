package models

import (
  "time"
)

type Formrequest struct {
	Formrequestid int `gorm:"column:formrequestid;primaryKey" json:"formrequestid"`
	Formreqtype int8 `gorm:"column:formreqtype" json:"formreqtype"`
	Formrequestno *string `gorm:"column:formrequestno" json:"formrequestno"`
	Formrequestdate time.Time `gorm:"column:formrequestdate" json:"formrequestdate"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Productplanid *int `gorm:"column:productplanid" json:"productplanid"`
	Slocfromid int `gorm:"column:slocfromid" json:"slocfromid"`
	Isjasa *int8 `gorm:"column:isjasa" json:"isjasa"`
	Requestedbyid int `gorm:"column:requestedbyid" json:"requestedbyid"`
	Description *string `gorm:"column:description" json:"description"`
	Useraccessid int `gorm:"column:useraccessid" json:"useraccessid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Productplan Productplan `gorm:"foreignKey:productplanid;references:productplanid"`
	Requestedby Requestedby `gorm:"foreignKey:requestedbyid;references:requestedbyid"`
	Sloc Sloc `gorm:"foreignKey:slocfromid;references:slocid"`
}

func (Formrequest) TableName() string {
	return "formrequest"
}
