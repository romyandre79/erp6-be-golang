package models

import "time"

type Chat struct {
	ChatID     int       `gorm:"primaryKey;column:chatid" json:"chatid"`
	SenderID   int       `gorm:"column:senderid" json:"sender_id"`
	ReceiverID int       `gorm:"column:receiverid" json:"receiver_id"`
	Message    string    `gorm:"column:message;type:text" json:"message"`
	IsRead     int       `gorm:"column:isread;default:0" json:"is_read"` // 0: Unread, 1: Read
	CreatedAt  time.Time `gorm:"column:createdat" json:"created_at"`
}

func (Chat) TableName() string {
	return "chat" // Or whatever convention usage
}
