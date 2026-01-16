package generator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"google.golang.org/protobuf/proto"
	waProto "go.mau.fi/whatsmeow/binary/proto"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"gorm.io/gorm"

	"erp6-be-golang/models"
)

var (
	waClients      = make(map[string]*whatsmeow.Client) // Key: Phone Number (e.g. "628...")
	waClientsMutex sync.RWMutex
	waContainer    *sqlstore.Container
	
	currentQRs     = make(map[string]string) // Key: Phone Number (target)
	currentQRLocks sync.Mutex
	
	waDB           *gorm.DB
)

// SetWADatabase sets the DB instance for AI processing
func SetWADatabase(db *gorm.DB) {
	waDB = db
}

// InitWhatmeow initializes the WhatsApp clients from SQL Store
func InitWhatmeow() error {
	waClientsMutex.Lock()
	defer waClientsMutex.Unlock()

	if waContainer != nil {
		return nil // Already initialized
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	
	// Ensure we have a DB connection
	if waDB == nil {
		return fmt.Errorf("database not initialized")
	}

	// GET RAW CONNECTION STRING
	// We need to construct the DSN for sqlstore "mysql"
	// Best way needs explicit config, but we can try to reuse the existing pool if supported, 
	// or just open a new one with same credentials. 
	// For now, let's assume standard MySQL DSN.
	// Since we can't easily extract DSN from GORM, we might need a workaround.
	// However, `sqlstore.New` takes a dialect and URI.
	// Let's assume we can get it from env or just use a generic DSN if previously set.
	// A better approach for integrated systems: Use the existing `*sql.DB`? 
	// sqlstore doesn't support passing *sql.DB directly in `New`, but `NewWithDB` might exist? 
	// Checking docs... `New` opens it. 
	// Workaround: We will use the standard DSN format.
	
	// Use SQLite for Whatsmeow as MySQL support is missing in this version
	// Store in a local file "wa-session.db"
	// Use SQLite for Whatsmeow as MySQL support is missing in this version
	// Store in a local file "wa-session.db"
	// modernc.org/sqlite: use _txlock=immediate to avoid nested transaction errors.
	// Enable WAL mode for better concurrency to fix "database is locked" errors.
	dsn := "file:wa-session.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_busy_timeout=30000&_txlock=immediate"
	container, err := sqlstore.New(context.Background(), "sqlite", dsn, dbLog)
	if err != nil {
		return fmt.Errorf("failed to connect to wa database: %v", err)
	}
	
	waContainer = container

	// Load all devices
	devices, err := container.GetAllDevices(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get devices: %v", err)
	}

	fmt.Printf("[WA] Found %d devices in database\n", len(devices))

	for _, device := range devices {
		// Identify phone number from JID
		if device.ID == nil {
			continue
		}
		
		jid := device.ID.ToNonAD()
		phone := jid.User
		
		fmt.Printf("[WA] Loading device: %s\n", phone)
		
		clientLog := waLog.Stdout("Client-"+phone, "DEBUG", true)
		client := whatsmeow.NewClient(device, clientLog)
		
		// Register Handler with Client Context
		client.AddEventHandler(func(evt interface{}) {
			eventHandler(client, phone, evt)
		})
		
		if err := client.Connect(); err != nil {
			fmt.Printf("[WA] Failed to connect %s: %v\n", phone, err)
		} else {
			fmt.Printf("[WA] Connected %s\n", phone)
		}
		
		waClients[phone] = client
		
		// Sync Status
		updateDeviceStatus(phone, jid.String(), "Connected")
	}

	return nil
}

func updateDeviceStatus(phone, jid, status string) {
	if waDB == nil { return }
	
	// Upsert into whatsappdevice table
	var dev models.Whatsappdevice
	err := waDB.Where("phonenumber = ?", phone).First(&dev).Error
	if err != nil {
		// Create
		dev = models.Whatsappdevice{
			Phonenumber: phone,
			Jid:         jid,
			Nickname:    "WA " + phone,
			Status:      status,
			Lastactive:  time.Now(),
		}
		waDB.Create(&dev)
	} else {
		// Update
		dev.Status = status
		dev.Lastactive = time.Now()
		if jid != "" {
			dev.Jid = jid
		}
		waDB.Save(&dev)
	}
}

// GetLoginQR generates a QR code for a specific phone number (or new one)
func GetLoginQR(targetPhone string) (string, error) {
	if waContainer == nil {
		if err := InitWhatmeow(); err != nil {
			return "", err
		}
	}
	
	waClientsMutex.Lock()
	client, exists := waClients[targetPhone]
	waClientsMutex.Unlock()
	
	// If client exists and is connected
	if exists && client.IsConnected() {
		return "Already logged in", nil
	}
	
	// If client doesn't exist, create new device
	if !exists {
		// Create new device in store
		device := waContainer.NewDevice()
		clientLog := waLog.Stdout("Client-New", "DEBUG", true)
		client = whatsmeow.NewClient(device, clientLog)
		// We don't know the phone yet until they scan!
		// But the frontend expects to "register" a specific number? 
		// Actually, standard flow: Scan QR -> We get Event -> We know user.
		// So we map the temporary client, show QR. Once logged in, we rename/move it in map?
		// Actually, `waClients` key is PHONE. 
		// When we start a NEW Login, we don't know the phone number yet.
		// We can perhaps handle this by returning the *channel* to the frontend
		// and once logged in, we register it properly.
	}

	// For simplicity: If re-login for known phone, use that device.
	// If NEW, use a new device.
	
	if client.Store.ID != nil {
		// Already has session, maybe disconnected?
		err := client.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect existing: %v", err)
		}
		if client.IsConnected() {
			return "Logged in", nil
		}
	}

	// Convert channel to QR string?
	// Whatsmeow doesn't give a static QR string immediately, it pushes to channel.
	
	// IMPORTANT: Whatsmeow GetQRChannel returns a channel that emits events.
	// We need to capture the first code.
	
	qrChan, _ := client.GetQRChannel(context.Background())
	if err := client.Connect(); err != nil {
		return "", err
	}
	
	// Wait for first QR
	select {
	case evt := <-qrChan:
		if evt.Event == "code" {
			// Update the map for this requested phone (even if we don't know real phone yet)
			// Ideally we return this code to FE.
			return evt.Code, nil
		} else {
			return "", fmt.Errorf("login event: %s", evt.Event)
		}
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timeout waiting for QR")
	}
}


// SendDocument sends a document (file) to a WhatsApp user
func SendDocument(jid, filePath, caption string) error {
	// Pick a default client or logic to select one?
	// For now, pick ANY connected client.
	waClientsMutex.RLock()
	var client *whatsmeow.Client
	for _, c := range waClients {
		if c.IsConnected() {
			client = c
			break
		}
	}
	waClientsMutex.RUnlock()
	
	if client == nil {
		return fmt.Errorf("no connected WhatsApp clients")
	}
	
	// ... (Rest of logic similar, just using `client`) ...
	
	// Parse JID
	recipient, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("invalid JID: %v", err)
	}

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Determine MIME type
	mimeType := "application/octet-stream"
	fileName := fileInfo.Name()
	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
		mimeType = "application/pdf"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		mimeType = "text/csv"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".png") {
		mimeType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".jpg") {
		mimeType = "image/jpeg"
	}

	// Upload file
	uploaded, err := client.Upload(context.Background(), fileData, whatsmeow.MediaDocument)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	// Create document message
	documentMsg := &waProto.DocumentMessage{
		URL:           proto.String(uploaded.URL),
		Mimetype:      proto.String(mimeType),
		Title:         proto.String(fileName),
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(fileData))),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		DirectPath:    proto.String(uploaded.DirectPath),
	}

	if caption != "" {
		documentMsg.Caption = proto.String(caption)
	}

	msg := &waProto.Message{
		DocumentMessage: documentMsg,
	}

	_, err = client.SendMessage(context.Background(), recipient, msg)
	return err
}

func authorizeWAUser(jid string) (*models.Useraccess, error) {
	if waDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	phone := strings.Split(jid, "@")[0]
	
	fmt.Printf("[WA Debug] Authorizing phone: '%s'\n", phone)
	
	var user models.Useraccess
	// Multi-device support does not change how we authorize External Users.
	// Users are still identified by THEIR phone number.
	if err := waDB.Where("CAST(phoneno AS CHAR) = ? OR CAST(userwa AS CHAR) = ?", phone, phone).First(&user).Error; err != nil {
		return nil, err
	}
	
	return &user, nil
}

// eventHandler now takes the specific client instance and its phone ID
func eventHandler(client *whatsmeow.Client, myPhone string, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}
		
		fmt.Printf("[WA %s] Received Message from %s\n", myPhone, v.Info.Sender.User)
		
		// Use the client that received the message for all operations
		
		actualJID := v.Info.Sender.ToNonAD()
		phoneNumber := actualJID.User
		if strings.Contains(actualJID.Server, "lid") {
			if v.Info.MessageSource.SenderAlt.User != "" {
				phoneNumber = v.Info.MessageSource.SenderAlt.User
			}
		}

		user, err := authorizeWAUser(phoneNumber)
		if err != nil {
			return
		}
		
		senderPhone := phoneNumber
		
		text := ""
		mediaType := ""
		var mediaData []byte
		var mediaErr error
		fileName := ""

		if v.Message.GetConversation() != "" {
			text = v.Message.GetConversation()
		} else if v.Message.GetExtendedTextMessage().GetText() != "" {
			text = v.Message.GetExtendedTextMessage().GetText()
		} else {
			// Check for Media
			if img := v.Message.GetImageMessage(); img != nil {
				mediaType = "image"
				text = img.GetCaption()
				mediaData, mediaErr = client.Download(context.Background(), v.Message.GetImageMessage())
				fileName = "image_" + v.Info.ID + ".jpg"
			} else if doc := v.Message.GetDocumentMessage(); doc != nil {
				mediaType = "document"
				text = doc.GetCaption()
				mediaData, mediaErr = client.Download(context.Background(), v.Message.GetDocumentMessage())
				fileName = doc.GetFileName()
			} // ... other types ...
		}

		if mediaType != "" {
			if mediaErr == nil {
				saveDir := fmt.Sprintf("./public/uploads/whatsapp/%s", time.Now().Format("2006-01-02"))
				os.MkdirAll(saveDir, 0755)
				savePath := filepath.Join(saveDir, fileName)
				os.WriteFile(savePath, mediaData, 0644)
				SendMessageViaClient(client, senderPhone, fmt.Sprintf("✅ File received: %s", fileName))
			}
			if text == "" { return }
		}
		
		// Create Request Context for Workflow
		if text != "" && waDB != nil {
			var reqCtx fasthttp.RequestCtx
			reqCtx.Request.Header.SetMethod("POST")
			reqCtx.Request.SetRequestURI("/whatsapp/aicommand")
			reqCtx.PostArgs().Set("command", text)
			
			app := fiber.New()
			c := app.AcquireCtx(&reqCtx)
			defer app.ReleaseCtx(c)
			
			c.Locals("userid", user.Useraccessid)
			c.Locals("username", user.Username)
			c.Locals("db", waDB)
			c.Locals("wfEngine", []WorkflowEngine{})

			// Callback using SPECIFIC CLIENT
			lastSentMsg := ""
			waCallback := func(msg string) {
				if msg != "" {
					recipient := senderPhone + "@s.whatsapp.net"
					// Send via the client that received the message
					if _, err := client.SendMessage(context.Background(), types.NewJID(senderPhone, types.DefaultUserServer), &waProto.Message{Conversation: proto.String(msg)}); err == nil {
						fmt.Printf("[WA %s] Real-time sent to %s: %s\n", myPhone, recipient, msg)
						lastSentMsg = msg
					} else {
						fmt.Printf("[WA %s] Failed send to %s: %v\n", myPhone, recipient, err)
					}
				}
			}
			
			c.Locals("wfExtras", map[string]interface{}{
				"send_wa_callback": waCallback,
			})
			
			params := map[string]interface{}{} 
			c.Locals("nestedWorkflow", false)
			
			if err := ExecuteFlow(c, waDB, "aicommand", false, params); err == nil {
				// Process results similar to before...
				// For brevity, using simplified result extraction
				// (You can copy full extraction logic if needed, but the Core concept is using client)
				
				if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok {
					msg := ""
					for _, node := range wfEngine {
						// ... logic to find msg ...
						if strings.EqualFold(node.ComponentName, "sendmessage") {
							// Check input
							if inputParams, ok := node.DataInputNode.(map[string]string); ok {
								if m, ok := inputParams["message"]; ok && m != "" {
									msg = ResolveParam(c, m)
								}
							}
						}
					}
					
					if msg != "" && msg != lastSentMsg {
						waCallback(msg)
					}
				}
			}
		}
	
	case *events.PairSuccess:
		// New Login!
		newID := v.ID.ToNonAD()
		newPhone := newID.User
		fmt.Printf("[WA] Pair Success! New device: %s\n", newPhone)
		
		waClientsMutex.Lock()
		waClients[newPhone] = client
		waClientsMutex.Unlock()
		
		// Update persistent map in DB
		updateDeviceStatus(newPhone, newID.String(), "Connected")
		
		// Re-register handler with correct phone
		client.RemoveEventHandlers()
		client.AddEventHandler(func(evt interface{}) {
			eventHandler(client, newPhone, evt)
		})
	}
}


func SendMessageViaClient(client *whatsmeow.Client, phone, text string) error {
	jid := types.NewJID(phone, types.DefaultUserServer)
	msg := &waProto.Message{Conversation: proto.String(text)}
	_, err := client.SendMessage(context.Background(), jid, msg)
	return err
}

func SendMessage(jidStr string, text string) error {
	// Pick ANY connected client
	waClientsMutex.RLock()
	var client *whatsmeow.Client
	for _, c := range waClients {
		if c.IsConnected() {
			client = c
			break
		}
	}
	waClientsMutex.RUnlock()
	
	if client == nil {
		return fmt.Errorf("no connected WhatsApp clients")
	}
	
	jid, err := types.ParseJID(jidStr)
	if err != nil {
		if !strings.Contains(jidStr, "@") {
			jidStr = jidStr + "@s.whatsapp.net"
			jid, err = types.ParseJID(jidStr)
		}
		if err != nil {
			return fmt.Errorf("invalid JID: %v", err)
		}
	}

	msg := &waProto.Message{
		Conversation: proto.String(text),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.SendMessage(ctx, jid, msg)
	return err
}


// SendWhatsAppMedia sends media from a URL
func SendWhatsAppMedia(jidStr, mediaURL, mediaType, caption string) error {
	// Pick ANY connected client
	waClientsMutex.RLock()
	var client *whatsmeow.Client
	for _, c := range waClients {
		if c.IsConnected() {
			client = c
			break
		}
	}
	waClientsMutex.RUnlock()
	
	if client == nil {
		return fmt.Errorf("no connected WhatsApp clients")
	}

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		if !strings.Contains(jidStr, "@") {
			jidStr = jidStr + "@s.whatsapp.net"
			jid, err = types.ParseJID(jidStr)
		}
		if err != nil {
			return fmt.Errorf("invalid JID: %v", err)
		}
	}

	// Download media
	resp, err := http.Get(mediaURL)
	if err != nil {
		return fmt.Errorf("failed to download media: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read media body: %v", err)
	}

	var msg *waProto.Message
	
	switch strings.ToLower(mediaType) {
	case "image":
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
		if err != nil {
			return err
		}
		msg = &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				Caption:       proto.String(caption),
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String("image/jpeg"), // Default
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			},
		}
	case "video":
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaVideo)
		if err != nil {
			return err
		}
		msg = &waProto.Message{
			VideoMessage: &waProto.VideoMessage{
				Caption:       proto.String(caption),
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String("video/mp4"), // Default
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			},
		}
	case "document":
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaDocument)
		if err != nil {
			return err
		}
		// Try to extract filename from URL
		filename := filepath.Base(mediaURL)
		if strings.Contains(filename, "?") {
			filename = strings.Split(filename, "?")[0]
		}
		
		msg = &waProto.Message{
			DocumentMessage: &waProto.DocumentMessage{
				Caption:       proto.String(caption),
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String("application/octet-stream"),
				Title:         proto.String(filename),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			},
		}
	default: 
		return fmt.Errorf("unsupported media type: %s", mediaType)
	}

	_, err = client.SendMessage(context.Background(), jid, msg)
	return err
}

// SendWhatsAppButtons sends a text message with buttons
func SendWhatsAppButtons(jidStr, text string, buttons []string) error {
	// Pick ANY connected client
	waClientsMutex.RLock()
	var client *whatsmeow.Client
	for _, c := range waClients {
		if c.IsConnected() {
			client = c
			break
		}
	}
	waClientsMutex.RUnlock()
	
	if client == nil {
		return fmt.Errorf("no connected WhatsApp clients")
	}

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		if !strings.Contains(jidStr, "@") {
			jidStr = jidStr + "@s.whatsapp.net"
			jid, err = types.ParseJID(jidStr)
		}
		if err != nil {
			return fmt.Errorf("invalid JID: %v", err)
		}
	}

	// Dynamic Buttons
	var waButtons []*waProto.ButtonsMessage_Button
	for i, btnText := range buttons {
		id := fmt.Sprintf("btn_%d", i+1)
		parts := strings.SplitN(btnText, ":", 2)
		if len(parts) == 2 {
			id = strings.TrimSpace(parts[0])
			btnText = strings.TrimSpace(parts[1])
		}
		
		waButtons = append(waButtons, &waProto.ButtonsMessage_Button{
			ButtonID: proto.String(id),
			ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(btnText),
			},
			Type: waProto.ButtonsMessage_Button_RESPONSE.Enum(),
		})
	}

	// NOTE: ButtonsMessage is deprecated/limited support on some devices.
	// Using TemplateMessage or InteractiveMessage is preferred but more complex.
	// We'll attempt basic ButtonsMessage first.
	msg := &waProto.Message{
		ButtonsMessage: &waProto.ButtonsMessage{
			ContentText: proto.String(text),
			Buttons:     waButtons,
			HeaderType:  waProto.ButtonsMessage_TEXT.Enum(),
		},
	}

	_, err = client.SendMessage(context.Background(), jid, msg)
	return err
}

func init() {
	RegisterComponent("Internal Whatsapp", func(ctx *WorkflowContext) error {
		var msg string
		var phoneNumbers []string
		var mediaURL string
		var mediaType string
		var btnStr string

		// 1. Get params
		for _, p := range ctx.Params {
			switch strings.ToLower(p.InputName) {
			case "message":
				msg = ResolveParam(ctx.FiberCtx, p.CompValue)
			case "phone_number", "phone", "sendto":
				val := ResolveParam(ctx.FiberCtx, p.CompValue)
				parts := strings.Split(val, ",")
				for _, part := range parts {
					clean := strings.TrimSpace(part)
					if clean != "" {
						phoneNumbers = append(phoneNumbers, clean)
					}
				}
			case "media_url", "image", "video", "url":
				mediaURL = ResolveParam(ctx.FiberCtx, p.CompValue)
			case "media_type", "type":
				mediaType = ResolveParam(ctx.FiberCtx, p.CompValue)
			case "buttons", "button":
				btnStr = ResolveParam(ctx.FiberCtx, p.CompValue)
			}
		}

		// 2. Fallback
		if msg == "" {
			if wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine); ok {
				for i := len(wfEngine) - 1; i >= 0; i-- {
					if res, ok := wfEngine[i].ResultNode.(map[string]interface{}); ok {
						if m, ok := res["message"].(string); ok && m != "" {
							msg = m
							break
						}
					}
				}
			}
		}

		if len(phoneNumbers) == 0 {
			return fmt.Errorf("phone_number is required")
		}

		fmt.Printf("[Internal Whatsapp] Broadcasting to %d recipients\n", len(phoneNumbers))

		for _, phone := range phoneNumbers {
			var err error
			
			// Priority: Buttons > Media > Text
			if btnStr != "" {
				buttons := strings.Split(btnStr, ",")
				// Trim spaces
				for i, b := range buttons {
					buttons[i] = strings.TrimSpace(b)
				}
				err = SendWhatsAppButtons(phone, msg, buttons)
			} else if mediaURL != "" {
				// Guess type if missing
				if mediaType == "" {
					if strings.HasSuffix(mediaURL, ".mp4") {
						mediaType = "video"
					} else if strings.HasSuffix(mediaURL, ".pdf") || strings.HasSuffix(mediaURL, ".doc") || strings.HasSuffix(mediaURL, ".docx") || strings.HasSuffix(mediaURL, ".xls") {
						mediaType = "document"
					} else {
						mediaType = "image"
					}
				}
				err = SendWhatsAppMedia(phone, mediaURL, mediaType, msg)
			} else {
				if msg == "" {
					fmt.Printf("[Internal Whatsapp] No message content to send to %s\n", phone)
					continue
				}
				err = SendMessage(phone, msg)
			}

			if err != nil {
				fmt.Printf("[Internal Whatsapp] Failed to send to %s: %v\n", phone, err)
			} else {
				fmt.Printf("[Internal Whatsapp] Sent to %s\n", phone)
			}
			time.Sleep(100 * time.Millisecond)
		}

		return nil
	})
}

// RegisterWARoutes registers the WhatsApp API routes
func RegisterWARoutes(app *fiber.App) {
	api := app.Group("/api/whatsapp")
	api.Get("/devices", HandleGetDevices)
	api.Get("/qr", HandleGetWAQR)
	api.Delete("/device/:phone", HandleDeleteDevice)
}

// HandleGetDevices returns the list of devices
func HandleGetDevices(c *fiber.Ctx) error {
	var devices []models.Whatsappdevice
	if waDB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not initialized"})
	}
	
	if err := waDB.Find(&devices).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Check connection status in memory map for accuracy
	waClientsMutex.RLock()
	for i, dev := range devices {
		if client, ok := waClients[dev.Phonenumber]; ok {
			if client.IsConnected() {
				devices[i].Status = "Connected"
			} else {
				devices[i].Status = "Disconnected"
			}
		} else {
			devices[i].Status = "Offline"
		}
	}
	waClientsMutex.RUnlock()
	
	return c.JSON(devices)
}

// HandleGetWAQR returns a generic QR code for a new device
func HandleGetWAQR(c *fiber.Ctx) error {
	qr, err := GetLoginQR("new") // Use "new" to always get a fresh QR/Client
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	if qr == "Already logged in" {
		return c.Status(400).JSON(fiber.Map{"error": "Already logged in (unexpected for 'new')"})
	}
	
	return c.JSON(fiber.Map{"qr": qr})
}

// HandleDeleteDevice logout and delete a device
func HandleDeleteDevice(c *fiber.Ctx) error {
	phone := c.Params("phone")
	if phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phone number required"})
	}
	
	waClientsMutex.Lock()
	client, ok := waClients[phone]
	if ok && client != nil {
		client.Logout(context.Background())
		client.Disconnect()
		delete(waClients, phone)
	}
	waClientsMutex.Unlock()
	
	// Remove from DB
	if waDB != nil {
		waDB.Where("phonenumber = ?", phone).Delete(&models.Whatsappdevice{})
	}
	
	return c.JSON(fiber.Map{"message": "Device deleted"})
}

