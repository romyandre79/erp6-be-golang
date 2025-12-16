package admin

import (
	"erp6-be-golang/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SaveChatMessage persists a chat message to the database
// SaveChatMessage persists a chat message to the database
func SaveChatMessage(db *gorm.DB, senderID, receiverID int, message, attachment string) error {
	chatLog := models.Chat{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Attachment: attachment,
		IsRead:     0,
		CreatedAt:  time.Now(),
	}
	return db.Create(&chatLog).Error
}

// GetChatHistoryHandler fetches chat history between the current user and a target user
func GetChatHistoryHandler(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userid").(int)
	targetID := c.QueryInt("target_id")

	if targetID == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Target ID required"})
	}

	var chats []models.Chat
	// Fetch messages where (sender=me AND receiver=target) OR (sender=target AND receiver=me)
	err := db.Where(
		db.Where("senderid = ? AND receiverid = ?", userID, targetID).
			Or("senderid = ? AND receiverid = ?", targetID, userID),
	).Order("createdat ASC").Find(&chats).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   chats,
	})
}
