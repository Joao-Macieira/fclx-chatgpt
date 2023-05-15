package entity

// Chat gpt configs, can get in openia documentation
type ChatConfig struct {
	Model            *Model
	Temperature      float32
	TopP             float32
	N                int
	Stop             []string
	MaxTokens        int
	PresencePenalty  float32
	FrequencyPenalty float32
}

type Chat struct {
	ID                   string
	UserID               string
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessage        []*Message
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}
