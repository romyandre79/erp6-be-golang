package models

import (
	"time"
)

type Splettertype struct {
	Splettertypeid int       `gorm:"column:splettertypeid;primaryKey" json:"splettertypeid"`
	Splettername   string    `gorm:"column:splettername" json:"splettername"`
	Description    string    `gorm:"column:description" json:"description"`
	Validperiod    int       `gorm:"column:validperiod" json:"validperiod"`
	Recordstatus   int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate     time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Splettertype) TableName() string {
	return "splettertype"
}
