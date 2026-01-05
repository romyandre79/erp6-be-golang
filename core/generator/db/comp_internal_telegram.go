package generator

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"os"
	"net/http"
	"io"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"erp6-be-golang/models"
)

var (
	tgBot      *tgbotapi.BotAPI
	tgBotMutex sync.Mutex
	tgDB       *gorm.DB
)

// Helper to check and mark update as processed using ATOMIC FILE LOCK
// This works even if Fiber is running in Prefork mode (multiple processes)
func isUpdateProcessed(updateID int) bool {
	// Ensure directory exists
	os.MkdirAll("./tmp/telegram_updates", 0755)
	
	lockFile := fmt.Sprintf("./tmp/telegram_updates/%d.lock", updateID)
	// Try to create a file EXCLUSIVELY.
	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	
	if err != nil {
		if os.IsExist(err) {
			// File exists = Already processed
			fmt.Printf("[Telegram Debug] DEDUPE HIT: Lock file exists for ID %d\n", updateID)
			return true
		}
		// Other error
		fmt.Printf("[Telegram Debug] DEDUPE ERROR for ID %d: %v. Assuming processed to be safe.\n", updateID, err)
		return true 
	}
	
	fmt.Printf("[Telegram Debug] DEDUPE PASS: Created lock file for ID %d\n", updateID)
	f.Close()
	return false
}

// SetTelegramDatabase sets the DB instance for Telegram processing
func SetTelegramDatabase(db *gorm.DB) {
	tgDB = db
}

// InitTelegram initializes the Telegram bot
func InitTelegram(token string) error {
	tgBotMutex.Lock()
	defer tgBotMutex.Unlock()

	if tgBot != nil {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %v", err)
	}

	tgBot = bot
	fmt.Printf("[Telegram] Authorized on account %s\n", bot.Self.UserName)

	return nil
}

func init() {
	fmt.Println("[Telegram] Registering component: Internal Telegram")
	RegisterComponent("Internal Telegram", func(ctx *WorkflowContext) error {
		// 1. Get Webhook Body from previous generic Webhook component
		var body map[string]interface{}
		
		// Look for "webhook" component result in the execution history (wfEngine)
		if wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine); ok {
			for i := len(wfEngine) - 1; i >= 0; i-- {
				node := wfEngine[i]
				if node.ComponentName == "webhook" {
					if resMap, ok := node.ResultNode.(map[string]interface{}); ok {
						if b, ok := resMap["body"].(map[string]interface{}); ok {
							body = b
						}
					}
					break
				}
			}
		}

		if body == nil {
			return fmt.Errorf("no webhook body found from previous nodes")
		}

		// 2. Parse Telegram Update
		importJSON, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %v", err)
		}

		var update tgbotapi.Update
		if err := json.Unmarshal(importJSON, &update); err != nil {
			return fmt.Errorf("failed to unmarshal update: %v", err)
		}

		fmt.Printf("[Telegram Debug] Received UpdateID: %d\n", update.UpdateID)

		// Strong Deduplication: ALWAYS use MessageID if available, as it is constant across retries.
		// UpdateID *might* change (though unlikely, better safe).
		dedupeID := 0
		if update.Message != nil {
			dedupeID = update.Message.MessageID
			fmt.Printf("[Telegram Debug] Using MessageID for Dedupe: %d\n", dedupeID)
		} else {
			dedupeID = update.UpdateID
			fmt.Printf("[Telegram Debug] No MessageID, using UpdateID: %d\n", dedupeID)
		}

		if dedupeID != 0 {
			if isUpdateProcessed(dedupeID) {
				fmt.Printf("[Telegram] Duplicate ID %d ignored (retry suppressed)\n", dedupeID)
				return nil
			}
		}

		if update.Message == nil {
			return nil // No message to process
		}

		message := update.Message
		chatID := message.Chat.ID
		text := message.Text
		telegramUserID := message.From.ID
		userName := message.From.UserName

		fmt.Printf("[Telegram] Processing message from %s (ID: %d): %s\n", userName, telegramUserID, text)
		
		// Handle Media/Files
		var fileID string
		var fileName string
		fileType := ""
		
		if message.Document != nil {
			fileID = message.Document.FileID
			fileName = message.Document.FileName
			fileType = "document"
		} else if message.Photo != nil && len(message.Photo) > 0 {
			// Photos are arrays of different sizes, get the largest one (last one)
			photo := message.Photo[len(message.Photo)-1]
			fileID = photo.FileID
			fileName = fmt.Sprintf("photo_%d.jpg", time.Now().UnixNano())
			fileType = "photo"
		} else if message.Video != nil {
			fileID = message.Video.FileID
			fileName = fmt.Sprintf("video_%d.mp4", time.Now().UnixNano())
			fileType = "video"
		} else if message.Audio != nil {
			fileID = message.Audio.FileID
			fileName = message.Audio.FileName
			if fileName == "" {
				fileName = fmt.Sprintf("audio_%d.mp3", time.Now().UnixNano())
			}
			fileType = "audio"
		} else if message.Voice != nil {
			fileID = message.Voice.FileID
			fileName = fmt.Sprintf("voice_%d.ogg", time.Now().UnixNano())
			fileType = "voice"
		}
		
		if fileID != "" {
			fmt.Printf("[Telegram] Received %s with FileID: %s\n", fileType, fileID)
			
			// Get File URL
			fileConfig := tgbotapi.FileConfig{FileID: fileID}
			fileInfo, err := tgBot.GetFile(fileConfig)
			if err != nil {
				fmt.Printf("[Telegram] Failed to get file info: %v\n", err)
				SendTelegramMessage(chatID, "❌ Failed to retrieve file info")
			} else {
				// Construct URL: https://api.telegram.org/file/bot<token>/<file_path>
				// tgbotapi helper:
				fileURL := fileInfo.Link(tgBot.Token)
				
				// Download File
				resp, err := http.Get(fileURL)
				if err != nil {
					fmt.Printf("[Telegram] Failed to download file: %v\n", err)
					SendTelegramMessage(chatID, "❌ Failed to download file")
				} else {
					defer resp.Body.Close()
					
					// Save File
					saveDir := fmt.Sprintf("./public/uploads/telegram/%s", time.Now().Format("2006-01-02"))
					os.MkdirAll(saveDir, 0755)
					
					// Sanitize filename lightly
					fileName = strings.ReplaceAll(fileName, "..", "")
					if fileName == "" { fileName = "unknown_file" }
					
					savePath := filepath.Join(saveDir, fileName)
					out, err := os.Create(savePath)
					if err != nil {
						fmt.Printf("[Telegram] Failed to create file: %v\n", err)
						SendTelegramMessage(chatID, "❌ Failed to save file on server")
					} else {
						_, err = io.Copy(out, resp.Body)
						out.Close()
						if err != nil {
							fmt.Printf("[Telegram] Failed to write file: %v\n", err)
							SendTelegramMessage(chatID, "❌ Failed to write file content")
						} else {
							fmt.Printf("[Telegram] File saved to: %s\n", savePath)
							SendTelegramMessage(chatID, fmt.Sprintf("✅ File received and saved: %s", fileName))
							
							// If text is empty (just caption?), use caption
							if text == "" && message.Caption != "" {
								text = message.Caption
								fmt.Printf("[Telegram] Using caption as command: %s\n", text)
							}
							
							// If still no text, we just saved the file. Return?
							if text == "" {
								return nil
							}
						}
					}
				}
			}
		}
		// 3. Authorize User
		user, err := authorizeTelegramUser(telegramUserID)
		if err != nil {
			fmt.Printf("[Telegram] Unauthorized: %v\n", err)
			SendTelegramMessage(chatID, "⛔ You are not registered. Please contact admin.")
			return fmt.Errorf("user not authorized: %v", err)
		}

		// 4. Get Conversation State (Using DB ID)
		stateJSON := getTelegramUserState(fmt.Sprintf("%d", user.Useraccessid))
		fmt.Printf("[Telegram] User %s (DB ID: %d) authorized. State loaded.\n", userName, user.Useraccessid)
		fmt.Printf("[Telegram Debug] State Content: %s\n", stateJSON); // DEBUG PRINT

		// 5. Setup for ExecuteFlow ("aicommand")
		
		// Inject Real-Time Callback (similar to WA)
		realMessageSent := false
		
		tgCallback := func(msg string) {
			if msg != "" {
				SendTelegramMessage(chatID, msg)
				fmt.Printf("[Telegram] Real-time message sent to %d: %s\n", chatID, msg)
				
				if !strings.Contains(msg, "Processing") {
					realMessageSent = true
					fmt.Printf("[Telegram Debug] Callback triggered (Msg: %s). realMessageSent=TRUE\n", msg)
				}
			}
		}

		ctx.FiberCtx.Locals("wfExtras", map[string]interface{}{
			"send_wa_callback": tgCallback,
		})
		// Also set directly to avoid map casting issues in sendmessage
		ctx.FiberCtx.Locals("send_wa_callback", tgCallback)

		ctx.FiberCtx.Locals("userid", user.Useraccessid)
		ctx.FiberCtx.Locals("username", user.Username)
		ctx.FiberCtx.Locals("db", tgDB)

        // Stop the Graph from proceeding to the next node (Double Execution Killer)
        // MOVED TO END: Setting it here causes child nodes in ExecuteFlow to consume/reset it!
        // ctx.FiberCtx.Locals("skipNavigation", true)
        fmt.Printf("[Telegram Debug] Manual Execution Mode.\n")
        
        // Pass command AND existing state to workflow
        params := map[string]interface{}{
            "command": text,
            "conversation_state": stateJSON, // Pass state so AI component can use it
        }

		// FORCE Root Workflow: Ensure we are not inheriting "nested" state
		ctx.FiberCtx.Locals("nestedWorkflow", false)

        if err := ExecuteFlow(ctx.FiberCtx, tgDB, "aicommand", false, params); err != nil {
            fmt.Printf("[Telegram] Workflow execution error: %v\n", err)
            SendTelegramMessage(chatID, fmt.Sprintf("❌ Error: %v", err))
            return nil
        }
       
        // 6. Process Results
        // Now valid because ExecuteFlow has populated wfEngine
        
		if wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine); ok {
            msg := ""
            excelFile := ""
            foundSendMessage := false
            
            fmt.Printf("[Telegram Debug] Starting Result Scan.\n")

            for i := len(wfEngine) - 1; i >= 0; i-- {
                // Debug node
                node := wfEngine[i]
                fmt.Printf("[Telegram Debug] Scanning Node %d: %s\n", i, node.ComponentName)
                
                resMap, ok := node.ResultNode.(map[string]interface{})
                if !ok {
                	continue
                }

                // SAVE STATE: Check for conversation_state update
                // Accept from ANY component (e.g. SaveLog might carry it too)
                // if strings.EqualFold(node.ComponentName, "AIAssistant") || ... { // Too restrictive!
                
                if newState, ok := resMap["conversation_state"].(string); ok && newState != "" {
                		// Only save if it looks like JSON or valid state
                		if strings.HasPrefix(newState, "{") {
							fmt.Printf("[Telegram] Saving conversation state from node '%s' for user %d\n", node.ComponentName, user.Useraccessid)
	 						fmt.Printf("[Telegram Debug] New State to Save: %s\n", newState)
	 						
	 						// Helper to save state
							conversationDir := "./tmp/ai_conversations"
							os.MkdirAll(conversationDir, 0755)
							conversationFile := fmt.Sprintf("%s/%d.json", conversationDir, user.Useraccessid)
							stateData := map[string]interface{}{
								"conversation_state": newState,
								"updated_at":         "now", // timestamp
							}
							if data, err := json.Marshal(stateData); err == nil {
								os.WriteFile(conversationFile, data, 0644)
							}
                		}
                }

                 // Check for file
                if file, ok := resMap["excel_file"].(string); ok && file != "" {
                    excelFile = file
                }

                // Check for message
				candidateMsg := ""
                if m, ok := resMap["message"].(string); ok && m != "" {
                    candidateMsg = m
                } else if m, ok := resMap["result"].(string); ok && m != "" {
                    candidateMsg = m
                } else if strings.EqualFold(node.ComponentName, "sendmessage") {
                     if inputParams, ok := node.DataInputNode.(map[string]string); ok {
                         if sendMsg, ok := inputParams["message"]; ok && sendMsg != "" {
                             candidateMsg = ResolveParam(ctx.FiberCtx, sendMsg)
                         }
                     }
                }

				if candidateMsg != "" {
                    if strings.Contains(candidateMsg, "Processing your request") {
                        continue
                    }
                    
                    candidateMsg = strings.TrimSpace(candidateMsg)
                    msg = candidateMsg
                    foundSendMessage = true
                }
                
                if foundSendMessage {
                    break
                }
            }

            // We have real-time callback, so we usually don't need to send 'msg' again
            // UNLESS the callback wasn't used (e.g. Scraper fallback).
            // But for AI, it uses SendMessage.
            // Let's rely on callback for main message.
            // Only send if NOT sent via callback? 
             
            if msg != "" && !realMessageSent {
                 SendTelegramMessage(chatID, msg)
                 fmt.Printf("[Telegram Debug] Sending fallback message (realMessageSent=FALSE): %s\n", msg)
            }

            if excelFile != "" {
                 SendTelegramMessage(chatID, "📂 File Generated: " + excelFile)
            }
            
            // Legacy/Fallback check: If no callback fired, send result?
            // Safe to ignore for now if AI workflow works.
        }
        
        // CRITICAL DE-DUPLICATION FIX:
        // Set skipNavigation to TRUE here, at the very end.
        // This ensures proper state saving and workflow execution happened above via manual call.
        // Now we tell the PARENT graph engine to STOP and not proceed to the "AI Command" node in the graph.
        // Doing this here avoids child nodes (in ExecuteFlow) from consuming/resetting the flag.
        fmt.Printf("[Telegram Debug] Stopping Graph Navigation to prevent Double Execution.\n")
        ctx.FiberCtx.Locals("skipNavigation", true)
        
		return nil
        
		return nil
	})
	
	// Optional: Register a "TelegramReply" component if they want explicit send
	RegisterComponent("TelegramReply", func(ctx *WorkflowContext) error {
		// Find message and chat_id
		var msg string
		var chatID int64
		
		// 1. Get message from input params (mapped in designer)
		for _, p := range ctx.Params {
			if strings.ToLower(p.InputName) == "message" {
				msg = ResolveParam(ctx.FiberCtx, p.CompValue)
			}
			if strings.ToLower(p.InputName) == "chat_id" {
				// Try to parse
				cidStr := ResolveParam(ctx.FiberCtx, p.CompValue)
				fmt.Sscanf(cidStr, "%d", &chatID)
			}
		}
		
		// 2. Fallback: Look in wfEngine for "Telegram" node (for chat_id) and last result (for message)
		if wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine); ok {
			// Find chat_id
			if chatID == 0 {
				for i := len(wfEngine) - 1; i >= 0; i-- {
					if wfEngine[i].ComponentName == "Telegram" {
						if res, ok := wfEngine[i].ResultNode.(map[string]interface{}); ok {
							if cid, ok := res["chat_id"].(int64); ok {
								chatID = cid
							}
						}
						break
					}
				}
			}
			
			// Find message (if not param)
			if msg == "" {
				// Look for last result with "message"
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
		
		if chatID != 0 && msg != "" {
			SendTelegramMessage(chatID, msg)
		}
		
		return nil
	})
}

func authorizeTelegramUser(telegramUserID int64) (*models.Useraccess, error) {
	if tgDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user models.Useraccess
	telegramIDStr := fmt.Sprintf("%d", telegramUserID)
	// Check usertelegram column
	if err := tgDB.Where("usertelegram = ?", telegramIDStr).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func SendTelegramMessage(chatID int64, text string) error {
	if tgBot == nil {
		return fmt.Errorf("Telegram bot not initialized")
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	_, err := tgBot.Send(msg)
	if err != nil {
		fmt.Printf("[Telegram] SendMessage Error: %v\n", err)
	}
	return err
}

// Conversation State Helpers
var tgUserStates = make(map[string]string)
var tgStateMutex sync.Mutex

func getTelegramUserState(userID string) string {
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

func saveTelegramUserState(userID, state string) {
	tgStateMutex.Lock()
	defer tgStateMutex.Unlock()
	tgUserStates[userID] = state
}
