package models

import (
	"time"
)

type Modulerelation struct {
	Modulerelationid int       `gorm:"column:modulerelationid;primaryKey" json:"modulerelationid"`
	Moduleid         int       `gorm:"column:moduleid" json:"moduleid"`
	Relationid       int       `gorm:"column:relationid" json:"relationid"`
	Updatedate       time.Time `gorm:"column:updatedate" json:"updatedate"`
	Modules          Modules   `gorm:"foreignKey:moduleid;references:moduleid"`
	Relation         Modules   `gorm:"foreignKey:relationid;references:moduleid"`
}

func (Modulerelation) TableName() string {
	return "modulerelation"
}
