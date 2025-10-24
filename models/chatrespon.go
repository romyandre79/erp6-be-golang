package models

import (
	"time"
)

type Chatrespon struct {
	Chatresponid int       `gorm:"column:chatresponid;primaryKey" json:"chatresponid"`
	Msgreply     string    `gorm:"column:msgreply" json:"msgreply"`
	Sourcedata   string    `gorm:"column:sourcedata" json:"sourcedata"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Chatrespon) TableName() string {
	return "chatrespon"
}
