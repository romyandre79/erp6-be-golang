package models

import (
	"time"
)

type Whatsappdevice struct {
	Whatsappdeviceid int       `gorm:"column:whatsappdeviceid;primaryKey;autoIncrement" json:"whatsappdeviceid"`
	Phonenumber      string    `gorm:"column:phonenumber;not null" json:"phonenumber"`
	Nickname         string    `gorm:"column:nickname;not null" json:"nickname"`
	Jid              string    `gorm:"column:jid;not null" json:"jid"`
	Status           string    `gorm:"column:status;not null" json:"status"`
	Lastactive       time.Time `gorm:"column:lastactive;not null;default:CURRENT_TIMESTAMP" json:"lastactive"`
}

func (Whatsappdevice) TableName() string {
	return "whatsappdevice"
}
