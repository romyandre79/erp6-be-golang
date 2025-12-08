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

// Message types
type WsPayload struct {
	Type     string          `json:"type"`
	TargetID int             `json:"target_id,omitempty"` // UserID to send to
	Data     json.RawMessage `json:"data"`
}

// GlobalDB is set by routes.go for WebSocket handlers to access
var GlobalDB *gorm.DB

func WebSocketHandler(c *websocket.Conn) {
	log.Println("WebSocketHandler: Started")
	useridVal := c.Locals("userid")
	if useridVal == nil {
		c.Close()
		return
	}
	useraccessid := useridVal.(int)

	info := &ws.RegisterInfo{UserID: useraccessid, Conn: c}
	ws.GlobalHub.Register <- info

	// Status updates handled by Hub callbacks

	defer func() {
		ws.GlobalHub.Unregister <- info
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("WS Read Error:", err)
			}
			break
		}

		log.Printf("WS Recv Raw: %s\n", string(msg))

		// Parse generic payload
		var payload WsPayload
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("WS Parse Error:", err)
			continue
		}

		// If TargetID is present, route it
		if payload.TargetID > 0 {
			// Save to Database if type is 'chat'
			if payload.Type == "chat" {
				// We need to extract the text from the data map.
				var dataMap map[string]interface{}
				if err := json.Unmarshal(payload.Data, &dataMap); err == nil {
					if text, ok := dataMap["text"].(string); ok {
						// Get DB instance from somewhere - we need to pass it
						// For now, use a global or inject it. Let's use a package-level var.
						if GlobalDB != nil {
							go SaveChatMessage(GlobalDB, useraccessid, payload.TargetID, text)
						}
					}
				}
			}

			// Let's re-wrap to ensure sender is known
			outData := map[string]interface{}{
				"type":      payload.Type,
				"sender_id": useraccessid,
				"data":      payload.Data,
			}
			outBytes, _ := json.Marshal(outData)

			log.Printf("WS Routing to %d: %s\n", payload.TargetID, string(outBytes))
			ws.GlobalHub.SendToUser(payload.TargetID, outBytes)
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

	ws.GlobalHub.SendToUser(req.ReceiverID, payload)

	return c.JSON(fiber.Map{"status": "success", "data": todo})
}
