package models

import (
	"time"
)

type Bookcatalog struct {
	Bookcatalogid   int       `gorm:"column:bookcatalogid;primaryKey" json:"bookcatalogid"`
	Bookcatalogcode string    `gorm:"column:bookcatalogcode" json:"bookcatalogcode"`
	Description     string    `gorm:"column:description" json:"description"`
	Parentid        int       `gorm:"column:parentid" json:"parentid"`
	Recordstatus    int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate      time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Bookcatalog) TableName() string {
	return "bookcatalog"
}
