package generator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	
	"google.golang.org/protobuf/proto"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	
	"gorm.io/gorm"
	_ "modernc.org/sqlite" 

	"erp6-be-golang/models"
)

var (
	waClient      *whatsmeow.Client
	waClientMutex sync.Mutex
	currentQR     string
	currentQRLock sync.Mutex
	waDB          *gorm.DB
)

// SetWADatabase sets the DB instance for AI processing
func SetWADatabase(db *gorm.DB) {
	waDB = db
}

// InitWhatmeow initializes the WhatsApp client
func InitWhatmeow() error {
	waClientMutex.Lock()
	defer waClientMutex.Unlock()

	if waClient != nil {
		if waClient.IsConnected() {
			return nil
		}
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	dbPath := "wa.db" 

	container, err := sqlstore.New(context.Background(), "sqlite", "file:"+dbPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbLog)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	waClient = whatsmeow.NewClient(deviceStore, clientLog)
	waClient.AddEventHandler(eventHandler)

	if waClient.Store.ID == nil {
		qrChan, _ := waClient.GetQRChannel(context.Background())
		
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					currentQRLock.Lock()
					currentQR = evt.Code
					currentQRLock.Unlock()
					
					fmt.Println("QR Code:", evt.Code) 
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}()

		err = waClient.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	} else {
		err = waClient.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}

	return nil
}

// SendDocument sends a document (file) to a WhatsApp user
func SendDocument(jid, filePath, caption string) error {
	if waClient == nil {
		return fmt.Errorf("WhatsApp client not initialized")
	}

	if !waClient.IsConnected() {
		return fmt.Errorf("WhatsApp client not connected")
	}

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

	// Determine MIME type based on extension
	mimeType := "application/octet-stream"
	fileName := fileInfo.Name()
	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
		mimeType = "application/pdf"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		mimeType = "text/csv"
	}

	// Upload file to WhatsApp servers
	uploaded, err := waClient.Upload(context.Background(), fileData, whatsmeow.MediaDocument)
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

	_, err = waClient.SendMessage(context.Background(), recipient, msg)
	return err
}

func authorizeWAUser(jid string) (*models.Useraccess, error) {
	if waDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	phone := strings.Split(jid, "@")[0]
	
	fmt.Printf("[WA Debug] Authorizing phone: '%s' (type: %T)\n", phone, phone)
	
	var user models.Useraccess
	// Ensure phone is treated as string by using CAST or explicit string comparison
	if err := waDB.Where("CAST(phoneno AS CHAR) = ? OR CAST(userwa AS CHAR) = ?", phone, phone).First(&user).Error; err != nil {
		return nil, err
	}
	
	return &user, nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}
		
		fmt.Printf("[WA Debug] Received Message from %s: ID %s\n", v.Info.Sender.User, v.Info.ID)

		// Debug logs for JID
		fmt.Printf("[WA Debug] Raw JID: %+v\n", v.Info.Sender)
		fmt.Printf("[WA Debug] JID.User: %s (type: %T)\n", v.Info.Sender.User, v.Info.Sender.User)
		fmt.Printf("[WA Debug] JID.String(): %s\n", v.Info.Sender.String())
		fmt.Printf("[WA Debug] MessageInfo: %+v\n", v.Info)
		
		// Get actual phone number - ToNonAD doesn't work if LID mapping not available
		// Use SenderAlt which contains the actual phone number
		actualJID := v.Info.Sender.ToNonAD()
		fmt.Printf("[WA Debug] Actual JID (ToNonAD): %s\n", actualJID.String())
		
		// If ToNonAD didn't work (still returns LID), use SenderAlt
		phoneNumber := actualJID.User
		if strings.Contains(actualJID.Server, "lid") {
			// Still a LID, extract from SenderAlt which has the real phone number
			if v.Info.MessageSource.SenderAlt.User != "" {
				phoneNumber = v.Info.MessageSource.SenderAlt.User
				fmt.Printf("[WA Debug] Extracted phone from SenderAlt: %s\n", phoneNumber)
			} else {
				fmt.Printf("[WA Debug] ToNonAD still returned LID and SenderAlt is empty\n")
			}
		}

		user, err := authorizeWAUser(phoneNumber)
		if err != nil {
			//fmt.Printf("[WA Debug] Unauthorized access from %s: %v\n", v.Info.Sender.User, err)
			//SendMessage(v.Info.Sender.User+"@s.whatsapp.net", "Please register your phone number in the system to use this service.")
			return
		}
		fmt.Printf("[WA Debug] Authorized User: %s (ID: %d)\n", user.Realname, user.Useraccessid)
		userIDStr := fmt.Sprintf("%d", user.Useraccessid)
		
		// Use actual phone number for replies
		senderPhone := phoneNumber

		text := ""
		if v.Message.GetConversation() != "" {
			text = v.Message.GetConversation()
		} else if v.Message.GetExtendedTextMessage().GetText() != "" {
			text = v.Message.GetExtendedTextMessage().GetText()
		}
		
		fmt.Printf("[WA Debug] extracted text: %s\n", text)

		if text != "" && waDB != nil {
			// Get conversation state for this user
			stateJSON := getWAUserState(senderPhone)
			dbDriver := waDB.Dialector.Name()
			result, err := processAI(text, stateJSON, dbDriver, userIDStr, waDB)
			if err != nil {
				fmt.Printf("[WA Debug] processAI Error: %v\n", err)
				SendMessage(senderPhone+"@s.whatsapp.net", fmt.Sprintf("Error: %v", err))
				return
			}
			
			fmt.Printf("[WA Debug] processAI Result: %+v\n", result)

			if newState, ok := result["conversation_state"].(string); ok {
				saveWAUserState(senderPhone, newState)
			}
			
			msg, _ := result["message"].(string)

			// Handle execution via comp_ai's ExecuteAIResult
			if exec, ok := result["execute"].(string); ok && exec == "true" {
				fmt.Printf("[WA Debug] Executing via ExecuteAIResult\n")
				
				// Create WorkflowContext
				wfCtx := &WorkflowContext{
					FiberCtx: nil,
					DB:       waDB,
					Params:   []WorkflowDetailResult{},
					Extras:   make(map[string]interface{}),
				}
				
				// Initialize wfEngine in Extras
				wfCtx.Extras["wfEngine"] = []WorkflowEngine{}
				
				// Call ExecuteAIResult from comp_ai
				if err := ExecuteAIResult(result, waDB, wfCtx); err != nil {
					fmt.Printf("[WA Debug] Execution error: %v\n", err)
					msg = fmt.Sprintf("❌ Execution Failed: %v", err)
				} else {
					// Debug: Print all wfEngine nodes
					if wfEngine, ok := wfCtx.Extras["wfEngine"].([]WorkflowEngine); ok {
						fmt.Printf("[WA Debug] wfEngine has %d nodes:\n", len(wfEngine))
						for i, node := range wfEngine {
							fmt.Printf("[WA Debug] Node %d: %s\n", i, node.ComponentName)
							fmt.Printf("[WA Debug]   Input: %+v\n", node.DataInputNode)
							fmt.Printf("[WA Debug]   Result: %+v\n", node.ResultNode)
						}
					}
					
					// Check if SendMessage component was executed (for formatted output)
					foundSendMessage := false
					if wfEngine, ok := wfCtx.Extras["wfEngine"].([]WorkflowEngine); ok {
						for _, node := range wfEngine {
							if node.ComponentName == "sendmessage" || node.ComponentName == "SendMessage" {
								foundSendMessage = true
								fmt.Printf("[WA Debug] Found SendMessage node\n")
								// Extract the message from SendMessage's input params
								if inputParams, ok := node.DataInputNode.(map[string]string); ok {
									fmt.Printf("[WA Debug] SendMessage input params: %+v\n", inputParams)
									if sendMsg, ok := inputParams["message"]; ok && sendMsg != "" {
										msg = sendMsg
										fmt.Printf("[WA Debug] Using formatted message from SendMessage: %s\n", sendMsg)
									} else if sendMsg, ok := inputParams["messagenotif"]; ok && sendMsg != "" {
										msg = sendMsg
										fmt.Printf("[WA Debug] Using formatted message from SendMessage (messagenotif): %s\n", sendMsg)
									}
								}
								break
							}
						}
					}
					
					if !foundSendMessage {
						// No SendMessage, show raw result
						msg += "\n\n✅ Command Executed Successfully."
						if execResult, ok := wfCtx.Extras["result"]; ok {
							msg += "\n\n📊 Result:\n"
							msg += formatResult(execResult)
						}
					}
				}
			}

			if msg != "" {
				// Check if there's an Excel file to send
				if excelFile, ok := result["excel_file"].(string); ok && excelFile != "" {
					fmt.Printf("[WA Debug] Sending Excel file: %s\n", excelFile)
					recipientJID := senderPhone + "@s.whatsapp.net"
					if err := SendDocument(recipientJID, excelFile, msg); err != nil {
						fmt.Printf("[WA Debug] SendDocument Error: %v\n", err)
						// Fallback to text message if file send fails
						if err := SendMessage(recipientJID, msg); err != nil {
							fmt.Printf("[WA Debug] SendMessage Error: %v\n", err)
						}
					} else {
						// File sent successfully, don't send text message
						fmt.Printf("[WA Debug] Excel file sent successfully\n")
					}
				} else {
					// No file, send text message
					recipientJID := senderPhone + "@s.whatsapp.net"
					fmt.Printf("[WA Debug] Sending reply to %s: %s\n", recipientJID, msg)
					if err := SendMessage(recipientJID, msg); err != nil {
						fmt.Printf("[WA Debug] SendMessage Error: %v\n", err)
					}
				}
			}
		}
	}
}

func formatResult(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		var lines []string
		for key, value := range v {
			if key == "final_url" && value == "" {
				continue // Skip empty final_url
			}
			lines = append(lines, fmt.Sprintf("  • %s: %v", key, value))
		}
		return strings.Join(lines, "\n")
	case []interface{}:
		if len(v) == 0 {
			return "  (empty)"
		}
		var lines []string
		for i, item := range v {
			if i >= 5 { // Limit to first 5 items
				lines = append(lines, fmt.Sprintf("  ... and %d more", len(v)-5))
				break
			}
			lines = append(lines, fmt.Sprintf("  %d. %v", i+1, item))
		}
		return strings.Join(lines, "\n")
	case string:
		if len(v) > 500 {
			return v[:500] + "..."
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Simple in-memory state store for WA users (should be DB in production)
var waUserStates = make(map[string]string)
var waStateMutex sync.Mutex

func getWAUserState(userID string) string {
	waStateMutex.Lock()
	defer waStateMutex.Unlock()
	return waUserStates[userID]
}

func saveWAUserState(userID, state string) {
	waStateMutex.Lock()
	defer waStateMutex.Unlock()
	waUserStates[userID] = state
}

func GetLoginQR() (string, error) {
	if waClient == nil {
		if err := InitWhatmeow(); err != nil {
			return "", err
		}
	}
	
	if waClient.IsConnected() && waClient.Store.ID != nil {
		return "Already logged in", nil
	}

	for i := 0; i < 60; i++ {
		currentQRLock.Lock()
		qr := currentQR
		currentQRLock.Unlock()

		if qr != "" {
			return qr, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return "", fmt.Errorf("timeout waiting for QR code generation")
}

func SendMessage(jidStr string, text string) error {
	if waClient == nil {
		return fmt.Errorf("client not initialized")
	}
	if !waClient.IsConnected() {
		return fmt.Errorf("client not connected")
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

	_, err = waClient.SendMessage(ctx, jid, msg)
	return err
}
