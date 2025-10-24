package models

import (
  "time"
)

type Appscompany struct {
	Appscompanyid int `gorm:"column:appscompanyid;primaryKey" json:"appscompanyid"`
	Appsid int `gorm:"column:appsid" json:"appsid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Appstoken string `gorm:"column:appstoken" json:"appstoken"`
	Enddate time.Time `gorm:"column:enddate" json:"enddate"`
	Rentamount float64 `gorm:"column:rentamount" json:"rentamount"`
	Recurringtypeid int `gorm:"column:recurringtypeid" json:"recurringtypeid"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Apps Apps `gorm:"foreignKey:appsid;references:appsid"`
}

func (Appscompany) TableName() string {
	return "appscompany"
}
