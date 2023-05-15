package entity

import "errors"

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
	ErasedMessages       []*Message
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}

func (chat *Chat) AddMessage(message *Message) error {
	if chat.Status == "ended" {
		return errors.New("Chat is ended, no more messages allowed")
	}

	for {
		if chat.Config.Model.GetMaxTokens() >= message.GetQuantityTokens()+chat.TokenUsage {
			chat.Messages = append(chat.Messages, message)
			chat.RefreshTokenUsage()
			break
		}

		chat.ErasedMessages = append(chat.ErasedMessages, chat.Messages[0])
		chat.Messages = chat.Messages[1:]
		chat.RefreshTokenUsage()
	}

	return nil
}

func (chat *Chat) RefreshTokenUsage() {
	chat.TokenUsage = 0

	for m := range chat.Messages {
		chat.TokenUsage += chat.Messages[m].GetQuantityTokens()
	}
}
