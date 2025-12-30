package generator

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"erp6-be-golang/models"
)

var (
	tgBot      *tgbotapi.BotAPI
	tgBotMutex sync.Mutex
	tgDB       *gorm.DB
)

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

		if update.Message == nil {
			return nil // No message to process
		}

		message := update.Message
		chatID := message.Chat.ID
		text := message.Text
		telegramUserID := message.From.ID
		userName := message.From.UserName

		fmt.Printf("[Telegram] Processing message from %s (ID: %d): %s\n", userName, telegramUserID, text)


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

		// 5. Prepare Result for Next Nodes (e.g. comp_ai)
		result := map[string]interface{}{
			"command":            text,
			"user_id":            user.Useraccessid, // Int
			"user_id_str":        fmt.Sprintf("%d", user.Useraccessid),
			"username":           user.Username,
			"conversation_state": stateJSON,
			"chat_id":            chatID,            // Needed for reply
			"telegram_user_id":   telegramUserID,
			"source":             "telegram",
		}

		// Store result for next node
		// Usually handled by framework returning this map, but explicit handling:
		// The `Execute` signature returns error. 
		// The framework appends the result to `wfEngine` using the return of this function?
		// Wait, `RegisterComponent` definition passes logic that returns error.
		// `ExecuteFlow` calls `component.Execute(ctx)`. 
		// `Execute` (in comp_ai.go example) appends to `wfEngine` manually?
		// Let's check `comp_ai.go`... yes, it does `ctx.FiberCtx.Locals("wfEngine", ...)`
		
		// So we must manually append result to wfEngine
		wm := WorkflowEngine{
			ComponentName: "Telegram",
			ResultNode:    result,
			Success:       true,
		}

		wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
		
		// ALSO set context params for `comp_ai` implicit lookup?
		// `comp_ai` looks at `ctx.Params` (explicit mapping) or `ctx.FiberCtx.Locals("userid")`.
		
		ctx.FiberCtx.Locals("userid", user.Useraccessid)
		ctx.FiberCtx.Locals("username", user.Username)
		ctx.FiberCtx.Locals("db", tgDB) // important for AI
		
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
