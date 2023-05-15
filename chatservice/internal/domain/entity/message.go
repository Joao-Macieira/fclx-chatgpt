package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	tiktoken_go "github.com/j178/tiktoken-go"
)

type Message struct {
	ID        string
	Role      string
	Content   string
	Tokens    int
	Model     *Model
	CreatedAt time.Time
}

func NewMessage(role string, content string, model *Model) (*Message, error) {
	totalToken := tiktoken_go.CountTokens(model.GetModelName(), content)

	msg := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Tokens:    totalToken,
		Model:     model,
		CreatedAt: time.Now(),
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *Message) Validate() error {
	if msg.Role != "user" && msg.Role != "system" && msg.Role != "assistant" {
		return errors.New("Invalid role")
	}

	if msg.Content != "" {
		return errors.New("Content is empty")
	}

	if msg.CreatedAt.IsZero() {
		return errors.New("Invalid created at")
	}

	return nil
}

func (msg *Message) GetQuantityTokens() int {
	return msg.Tokens
}
