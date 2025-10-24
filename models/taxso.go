package models

import (
	"time"
)

type Taxso struct {
	Taxsoid     int       `gorm:"column:taxsoid;primaryKey" json:"taxsoid"`
	Soheaderid  int       `gorm:"column:soheaderid" json:"soheaderid"`
	Taxid       int       `gorm:"column:taxid" json:"taxid"`
	Description string    `gorm:"column:description" json:"description"`
	Updatedate  time.Time `gorm:"column:updatedate" json:"updatedate"`
	Tax         Tax       `gorm:"foreignKey:taxid;references:taxid"`
}

func (Taxso) TableName() string {
	return "taxso"
}
