package models

import (
	"time"
)

type Repneraca struct {
	Repneracaid  int       `gorm:"column:repneracaid;primaryKey" json:"repneracaid"`
	Companyid    int       `gorm:"column:companyid" json:"companyid"`
	Nourut       int8      `gorm:"column:nourut" json:"nourut"`
	Accountcoded string    `gorm:"column:accountcoded" json:"accountcoded"`
	Keterangand  string    `gorm:"column:keterangand" json:"keterangand"`
	Iscountd     int       `gorm:"column:iscountd" json:"iscountd"`
	Isboldd      int       `gorm:"column:isboldd" json:"isboldd"`
	Counttyped   int       `gorm:"column:counttyped" json:"counttyped"`
	Accountcodek string    `gorm:"column:accountcodek" json:"accountcodek"`
	Keterangank  string    `gorm:"column:keterangank" json:"keterangank"`
	Iscountk     int       `gorm:"column:iscountk" json:"iscountk"`
	Isboldk      int       `gorm:"column:isboldk" json:"isboldk"`
	Counttypek   int       `gorm:"column:counttypek" json:"counttypek"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company      Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Repneraca) TableName() string {
	return "repneraca"
}
