package models

import (
  "time"
)

type Bookmembertype struct {
	Bookmembertypeid int `gorm:"column:bookmembertypeid;primaryKey" json:"bookmembertypeid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Membertypename string `gorm:"column:membertypename" json:"membertypename"`
	Loanlimit int8 `gorm:"column:loanlimit" json:"loanlimit"`
	Memberperiod int `gorm:"column:memberperiod" json:"memberperiod"`
	Regisprice int `gorm:"column:regisprice" json:"regisprice"`
	Denda int `gorm:"column:denda" json:"denda"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Bookmembertype) TableName() string {
	return "bookmembertype"
}
