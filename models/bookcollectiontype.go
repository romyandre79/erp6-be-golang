package models

import (
  "time"
)

type Bookcollectiontype struct {
	Bookcollectiontypeid int `gorm:"column:bookcollectiontypeid;primaryKey" json:"bookcollectiontypeid"`
	Collectiontypename string `gorm:"column:collectiontypename" json:"collectiontypename"`
	Collectiontypecode string `gorm:"column:collectiontypecode" json:"collectiontypecode"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookcollectiontype) TableName() string {
	return "bookcollectiontype"
}
