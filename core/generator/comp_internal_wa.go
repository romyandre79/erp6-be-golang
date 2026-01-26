package generator

import (
	"context"
	"database/sql"
	"encoding/json"
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

	// Raw Connection
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

// Log message to database
func logWAMessage(direction, sender, recipient, msgType, content, caption, status string) {
	if waDB == nil { return }
	
	// Ensure table exists (simple lazy migration)
	if !waDB.Migrator().HasTable(&models.WhatsappMessage{}) {
		waDB.AutoMigrate(&models.WhatsappMessage{})
	}

	log := models.WhatsappMessage{
		Direction: direction,
		Sender:    sender,
		Recipient: recipient,
		Type:      msgType,
		Content:   content,
		Caption:   caption,
		Status:    status,
		CreatedAt: time.Now(),
	}
	waDB.Create(&log)
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

			senderPhone := phoneNumber
			
			content  := ""
			caption  := ""
			mediaType := "text"
			var mediaData []byte
			var mediaErr error
			fileName := ""

			if v.Message.GetConversation() != "" {
				content = v.Message.GetConversation()
			} else if v.Message.GetExtendedTextMessage().GetText() != "" {
				content = v.Message.GetExtendedTextMessage().GetText()
			} else if img := v.Message.GetImageMessage(); img != nil {
				mediaType = "image"
				caption = img.GetCaption()
				mediaData, mediaErr = client.Download(context.Background(), img)
				fileName = fmt.Sprintf("image_%s_%s.jpg", v.Info.ID, time.Now().Format("150405"))
			} else if doc := v.Message.GetDocumentMessage(); doc != nil {
				mediaType = "document"
				caption = doc.GetCaption()
				mediaData, mediaErr = client.Download(context.Background(), doc)
				fileName = doc.GetFileName()
				if fileName == "" { fileName = fmt.Sprintf("doc_%s", v.Info.ID) }
			} else {
				content = ""
				mediaType = "unknown"
			}

			// Save media file if exists (and not text)
			if mediaType != "" && mediaType != "text" {
				if mediaErr == nil {
					dateDir := time.Now().Format("2006-01-02")
					saveDir := fmt.Sprintf("./public/uploads/whatsapp/%s", dateDir)
					os.MkdirAll(saveDir, 0755)
					savePath := filepath.Join(saveDir, fileName)
					
					// Save File
					if err := os.WriteFile(savePath, mediaData, 0644); err == nil {
						// Set Content to Web Path for frontend display
						webPath := fmt.Sprintf("/uploads/whatsapp/%s/%s", dateDir, fileName)
						content = webPath
						
						// Image specific handling (if needed in future, currently just setting content)
						if med := v.Message.GetImageMessage(); med != nil {
							// If image had specific logic, place here. 
							// Previously it double-saved. Removed.
						} 
					} else {
						fmt.Printf("[WA Error] Failed to write file: %v\n", err)
						content = "[Media Write Error]"
					}
				} else {
					content = "[Media Download Failed]"
				}
			}

			// Fallback for empty text (e.g. deleted message, status update, etc)
			if content == "" && mediaType == "text" {
				content = "[Empty Message / Unsupported Type]"
			}

				// LOG INCOMING
				logWAMessage("incoming", senderPhone, myPhone, mediaType, content, caption, "received")

				user, err := authorizeWAUser(phoneNumber)
				// ... existing auth logic ...
				if err != nil {
					return
				}
				
				// Create Request Context for Workflow
				// For workflow, we usually process the command text. 
				// If it's an image with caption, 'caption' is the command.
				// If it's text, 'content' is the command.
				cmdText := content
				if mediaType != "text" {
					cmdText = caption
					// If media exists but no caption, force a trigger so AI receives the file context
					// usage: use valid file path so CompDocument can read it if mapped to 'file'
					if cmdText == "" && mediaType != "" {
						cmdText = fmt.Sprintf("./public%s", content)
					}
				}

				if cmdText != "" && waDB != nil {
					var reqCtx fasthttp.RequestCtx
					reqCtx.Request.Header.SetMethod("POST")
					reqCtx.Request.SetRequestURI("/whatsapp/aicommand")
					reqCtx.PostArgs().Set("command", cmdText)
					
					app := fiber.New()
					cCtx := app.AcquireCtx(&reqCtx)
					defer app.ReleaseCtx(cCtx)
					
					// Setup Locals
					cCtx.Locals("userid", user.Useraccessid)
					cCtx.Locals("username", user.Username)
					cCtx.Locals("db", waDB)
					cCtx.Locals("wfEngine", []WorkflowEngine{})

					// Callback using SPECIFIC CLIENT
					lastSentMsg := ""
					waCallback := func(msg string) {
						if msg != "" {
							recipient := senderPhone + "@s.whatsapp.net"
							// Send via the client that received the message
							if _, err := client.SendMessage(context.Background(), types.NewJID(senderPhone, types.DefaultUserServer), &waProto.Message{Conversation: proto.String(msg)}); err == nil {
								fmt.Printf("[WA %s] Real-time sent to %s: %s\n", myPhone, recipient, msg)
								logWAMessage("outgoing", myPhone, senderPhone, "text", msg, "", "sent")
								lastSentMsg = msg
							} else {
								fmt.Printf("[WA %s] Failed send to %s: %v\n", myPhone, recipient, err)
								logWAMessage("outgoing", myPhone, senderPhone, "text", msg, "", "failed")
							}
						}
					}
					
					cCtx.Locals("wfExtras", map[string]interface{}{
						"send_wa_callback": waCallback,
					})
					
					params := map[string]interface{}{} 
					
					// 1. Get Conversation State
					stateJSON := getWAUserState(fmt.Sprintf("%d", user.Useraccessid))
					params["conversation_state"] = stateJSON

					// 2. Pass File Paths if media
					if mediaType != "text" && mediaType != "" {
						// content contains the web path (e.g. /uploads/whatsapp/...)
						// We need absolute path or relative to root? 
						// comp_ai.go expects file paths.
						// Let's pass the web path, assuming comp_ai can handle it or we convert to local.
						// Actually comp_ai uses it to read file content potentially. 
						// But for now, let's pass what we have.
						// Better: Pass the LOCAL path we just wrote.
						// But we constructed local path inside the media block above. 
						// Reconstruct it or assume 'content' is fine? 
						// 'content' is "/uploads/whatsapp/..." (Web Path)
						// Local path is "./public" + content.
						localPath := fmt.Sprintf("./public%s", content)
						params["file_paths"] = localPath
						fmt.Printf("[WA] Sending file to AI: %s\n", localPath)
					}

					cCtx.Locals("nestedWorkflow", false)
					
					if err := ExecuteFlow(cCtx, waDB, "aicommand", false, params); err == nil {
						if wfEngine, ok := cCtx.Locals("wfEngine").([]WorkflowEngine); ok {
							msg := ""
							
							// Process Results & Save State
							for _, node := range wfEngine {
								// Check for state update
								if resMap, ok := node.ResultNode.(map[string]interface{}); ok {
									if newState, ok := resMap["conversation_state"].(string); ok && newState != "" {
										if strings.HasPrefix(newState, "{") {
											// Save to file
											conversationDir := "./tmp/ai_conversations"
											os.MkdirAll(conversationDir, 0755)
											conversationFile := fmt.Sprintf("%s/%d.json", conversationDir, user.Useraccessid)
											stateData := map[string]interface{}{
												"conversation_state": newState,
												"updated_at":         time.Now().Format(time.RFC3339),
											}
											if data, err := json.Marshal(stateData); err == nil {
												os.WriteFile(conversationFile, data, 0644)
												fmt.Printf("[WA] Saved conversation state for user %d\n", user.Useraccessid)
											}
										}
									}
								}
								
								if strings.EqualFold(node.ComponentName, "sendmessage") {
									if inputParams, ok := node.DataInputNode.(map[string]string); ok {
										if m, ok := inputParams["message"]; ok && m != "" {
											msg = ResolveParam(cCtx, m)
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
	
	status := "sent"
	if err != nil { status = "failed" }
	// We don't easily know "myPhone" here without reverse lookup or storing it in client struct wrapper.
	// But we can approximate or leave sender empty.
	// Actually, client.Store.ID contains our JID.
	myJID := client.Store.ID.ToNonAD()
	
	logWAMessage("outgoing", myJID.User, jid.User, "text", text, "", status)
	
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
	api.Get("/sessions", HandleGetWASessions)
	api.Get("/logs", HandleGetWALogs)
	api.Get("/qr", HandleGetWAQR)
	api.Delete("/device/:phone", HandleDeleteDevice)
}

// HandleGetWALogs returns chat logs
func HandleGetWALogs(c *fiber.Ctx) error {
	if waDB == nil { return c.Status(500).JSON(fiber.Map{"error": "DB not init"}) }
	
	var logs []models.WhatsappMessage
	// Pagination? Limit to last 100 for now
	// Use Debug() to see the generated SQL
	result := waDB.Debug().Order("created_at desc").Limit(100).Find(&logs)
	if result.Error != nil {
		fmt.Printf("[WA Log API Error] %v\n", result.Error)
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	
	fmt.Printf("[WA Log API] Query found %d logs\n", len(logs))
	if len(logs) > 0 {
		fmt.Printf("[WA Log API] Sample - Content: '%s', Type: '%s'\n", logs[0].Content, logs[0].Type)
	}
	
	return c.JSON(logs)
}

// HandleGetWASessions returns raw data from the SQLite session file
func HandleGetWASessions(c *fiber.Ctx) error {
	// Open the SQLite file in read-only mode, share cache
	db, err := sql.Open("sqlite", "file:wa-session.db?mode=ro")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("failed to open sqlite: %v", err)})
	}
	defer db.Close()
	
	// Query specific columns that definitely exist (exclude adv_account which caused error)
	rows, err := db.Query("SELECT jid, registration_id, platform FROM whatsmeow_device")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("query error: %v", err)})
	}
	defer rows.Close()
	
	var results []map[string]interface{}
	
	for rows.Next() {
		var jid string
		var regID int
		var platform sql.NullString
		
		if err := rows.Scan(&jid, &regID, &platform); err != nil {
			continue
		}
		
		results = append(results, map[string]interface{}{
			"jid": jid,
			"registration_id": regID,
			"platform": platform.String,
			"has_adv": false, // Placeholder as column is missing
		})
	}
	
	return c.JSON(results)
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


// Conversation State Helpers
var waUserStates = make(map[string]string)
var waStateMutex sync.Mutex

func getWAUserState(userID string) string {
	// Read from file ./tmp/ai_conversations/{userID}.json
	path := fmt.Sprintf("./tmp/ai_conversations/%s.json", userID)
	data, err := os.ReadFile(path)
	if err != nil {
		return "{}"
	}
	
	var stateMap map[string]interface{}
	if err := json.Unmarshal(data, &stateMap); err != nil {
		return "{}"
	}
	
	if state, ok := stateMap["conversation_state"].(string); ok {
		return state
	}
	// Fallback/Legacy
	return string(data)
}

func saveWAUserState(userID, state string) {
	waStateMutex.Lock()
	defer waStateMutex.Unlock()
	waUserStates[userID] = state
}
