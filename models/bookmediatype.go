package models

import (
  "time"
)

type Bookmediatype struct {
	Bookmediatypeid int `gorm:"column:bookmediatypeid;primaryKey" json:"bookmediatypeid"`
	Mediatypename string `gorm:"column:mediatypename" json:"mediatypename"`
	Mediatypecode string `gorm:"column:mediatypecode" json:"mediatypecode"`
	Mediatypeicon string `gorm:"column:mediatypeicon" json:"mediatypeicon"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookmediatype) TableName() string {
	return "bookmediatype"
}
