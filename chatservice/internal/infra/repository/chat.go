package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/joao-macieira/fclx/chatservice/internal/domain/entity"
	"github.com/joao-macieira/fclx/chatservice/internal/infra/db"
)

type ChatRepositoryMySQL struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySQL(database *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		DB:      database,
		Queries: db.New(database),
	}
}

func (repository *ChatRepositoryMySQL) CreateChat(ctx context.Context, chat *entity.Chat) error {
	err := repository.Queries.CreateChat(
		ctx,
		db.CreateChatParams{
			ID:               chat.ID,
			UserID:           chat.UserID,
			InitialMessageID: chat.InitialSystemMessage.Content,
			Status:           chat.Status,
			TokenUsage:       int32(chat.TokenUsage),
			Model:            chat.Config.Model.Name,
			ModelMaxTokens:   int32(chat.Config.Model.MaxToken),
			Temperature:      float64(chat.Config.Temperature),
			TopP:             float64(chat.Config.TopP),
			N:                int32(chat.Config.N),
			Stop:             chat.Config.Stop[0],
			MaxTokens:        int32(chat.Config.MaxTokens),
			PresencePenalty:  float64(chat.Config.PresencePenalty),
			FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	)

	if err != nil {
		return err
	}

	err = repository.Queries.AddMessage(
		ctx,
		db.AddMessageParams{
			ID:        chat.InitialSystemMessage.ID,
			ChatID:    chat.ID,
			Content:   chat.InitialSystemMessage.Content,
			Role:      chat.InitialSystemMessage.Role,
			Tokens:    int32(chat.InitialSystemMessage.Tokens),
			CreatedAt: chat.InitialSystemMessage.CreatedAt,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *ChatRepositoryMySQL) FindByChatID(ctx context.Context, chatID string) (*entity.Chat, error) {
	chat := &entity.Chat{}

	response, err := repository.Queries.FindChatByID(ctx, chatID)

	if err != nil {
		return nil, errors.New("chat not found")
	}

	chat.ID = response.ID
	chat.UserID = response.UserID
	chat.Status = response.Status
	chat.TokenUsage = int(response.TokenUsage)
	chat.Config = &entity.ChatConfig{
		Model: &entity.Model{
			Name:     response.Model,
			MaxToken: int(response.ModelMaxTokens),
		},
		Temperature:      float32(response.Temperature),
		TopP:             float32(response.TopP),
		N:                int(response.N),
		Stop:             []string{response.Stop},
		MaxTokens:        int(response.ModelMaxTokens),
		PresencePenalty:  float32(response.FrequencyPenalty),
		FrequencyPenalty: float32(response.FrequencyPenalty),
	}

	chatMessages, err := repository.Queries.FindMessagesByChatID(ctx, chat.ID)

	if err != nil {
		return nil, err
	}

	for _, message := range chatMessages {
		chat.Messages = append(chat.Messages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	chatErasedMessages, err := repository.Queries.FindErasedMessagesByChatID(ctx, chat.ID)

	if err != nil {
		return nil, err
	}

	for _, message := range chatErasedMessages {
		chat.ErasedMessages = append(chat.Messages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	return chat, nil
}
