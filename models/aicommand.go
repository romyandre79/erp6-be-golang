package models

// AICommand represents the database structure for AI commands
type AICommand struct {
	AICommandID    int    `gorm:"primaryKey;column:aicommandid"`
	Name           string `gorm:"column:name"`
	Description    string `gorm:"column:description"`
	Questions      string `gorm:"column:questions"`
	Queries        string `gorm:"column:queries"`
	Params         string `gorm:"column:params"`
	Defaults       string `gorm:"column:defaults"`
	MetaAction     string `gorm:"column:metaaction"`
	SuccessMessage string `gorm:"column:successmessage"`
	Triggers       string `gorm:"column:triggers"`
}

// TableName overrides the table name used by User to `aicommand`
func (AICommand) TableName() string {
	return "aicommand"
}

// AIConversationState represents the current state of a conversation
type AIConversationState struct {
	EntityType          string            `json:"entity_type"`          // module, menu, table, workflow
	CurrentStep         int               `json:"current_step"`         // Current question index
	CollectedData       map[string]string `json:"collected_data"`       // Data collected so far
	IsComplete          bool              `json:"is_complete"`          // Whether conversation is complete
	WaitingConfirmation bool              `json:"waiting_confirmation"` // Waiting for user to type 'execute'
	History             []AIMessage       `json:"history"`              // Full conversation history
}

type AIMessage struct {
	Role      string `json:"role"` // "user", "assistant", "system"
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// AIQuestionFlow defines the questions for each entity type and query templates
type AIQuestionFlow struct {
	Questions      []AIQuestion      `json:"questions"`
	Queries        map[string]string `json:"queries"`
	ListQueries    map[string]string `json:"list_queries"`
	Params         []string          `json:"params"`
	Defaults       map[string]string `json:"defaults"`
	SuccessMessage string            `json:"success_message"`
	MetaAction     string            `json:"meta_action"`
	Description    string            `json:"description"`
	Triggers       []string          `json:"triggers"`
}

type AIQuestion struct {
	Key          string   `json:"key"`           // Field name to store answer
	Text         string   `json:"text"`          // Question to ask user
	Validation   string   `json:"validation"`    // Validation type: required, optional, number, etc.
	Options      []string `json:"options"`       // For select-type questions
	OptionsQuery string   `json:"options_query"` // SQL query to fetch options dynamically
	Description  string   `json:"description"`   // Help text
}

// EntityInfo holds information about an available command entity
type AIEntityInfo struct {
	Name        string
	Description string
	Triggers    string
}

// AIConfig holds configuration for the AI service
type AIConfig struct {
	Token    string
	Provider string
	Model    string
	BaseURL  string // Custom base URL for providers like Ollama
}
