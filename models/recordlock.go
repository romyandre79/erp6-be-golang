package models

import "time"

// Recordlock represents a lock on a database record to prevent concurrent editing
type Recordlock struct {
	Recordlockid int       `gorm:"primaryKey;autoIncrement" json:"recordlockid"`
	Tablename    string    `gorm:"type:varchar(100);not null" json:"tablename"`
	Recordid     int       `gorm:"not null" json:"recordid"`
	Lockedby     int       `gorm:"not null" json:"lockedby"`
	Lockedat     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"lockedat"`
	Locktype     string    `gorm:"type:varchar(50);default:'edit'" json:"locktype"`
	Sessionid    string    `gorm:"type:varchar(255)" json:"sessionid"`
}

func (Recordlock) TableName() string {
	return "recordlock"
}
