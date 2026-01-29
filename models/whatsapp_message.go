package models

import (
	"time"

	"gorm.io/gorm"
)

type WhatsappMessage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Direction string    `gorm:"size:20;index" json:"direction"` // "incoming" or "outgoing"
	Sender    string    `gorm:"size:50;index" json:"sender"`
	Recipient string    `gorm:"size:50;index" json:"recipient"`
	Type      string    `gorm:"size:20" json:"type"` // "text", "image", "document", etc.
	Content   string    `gorm:"type:text" json:"content"`
	Caption   string    `gorm:"type:text" json:"caption"` 
	Status    string    `gorm:"size:20" json:"status"` // "sent", "failed", "received"
	Raw       string    `gorm:"type:text" json:"-"`    // Optional: raw debug usage
}

func (WhatsappMessage) TableName() string {
	return "whatsappmessage"
}
