package admin

import (
	"encoding/json"
	"erp6-be-golang/core/ws"
	"erp6-be-golang/models"

	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var Hub *ws.Hub

// WebSocketHandler handles the websocket connection
func WebSocketHandler(c *websocket.Conn) {
	// useraccessid is set in locals by middleware
	useraccessid := c.Locals("userid").(int)

	log.Printf("User %d connected\n", useraccessid)

	info := &ws.RegisterInfo{
		UserID: useraccessid,
		Conn:   c,
	}

	Hub.Register <- info

	defer func() {
		Hub.Unregister <- info
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func GetUnreadNotifications(c *fiber.Ctx, db *gorm.DB) error {
	useraccessid := c.Locals("userid").(int)

	var todos []models.Usertodo
	result := db.Where("useraccessid = ? AND isread = 0", useraccessid).
		Order("tododate desc").
		Find(&todos)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   todos,
	})
}

func MarkAsRead(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	useraccessid := c.Locals("userid").(int)

	result := db.Model(&models.Usertodo{}).
		Where("usertodoid = ? AND useraccessid = ?", id, useraccessid).
		Update("isread", 1)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

type SendNotificationRequest struct {
	ReceiverID int    `json:"receiverid"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	DocNo      string `json:"docno"`
}

func SendNotification(c *fiber.Ctx, db *gorm.DB) error {
	var req SendNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	todo := models.Usertodo{
		Useraccessid: req.ReceiverID,
		Menuname:     req.Title,
		Description:  req.Message,
		Docno:        req.DocNo,
		Isread:       0,
	}

	if err := db.Create(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Broadcast via WS
	payload, _ := json.Marshal(map[string]interface{}{
		"type": "notification",
		"data": todo,
	})

	Hub.SendToUser(req.ReceiverID, payload)

	return c.JSON(fiber.Map{"status": "success", "data": todo})
}
