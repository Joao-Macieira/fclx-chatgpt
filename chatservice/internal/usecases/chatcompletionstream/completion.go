package chatcompletionstream

import (
	"context"
	"errors"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/joao-macieira/fclx/chatservice/internal/domain/entity"
	"github.com/joao-macieira/fclx/chatservice/internal/domain/gateway"
)

type ChatCompletionConfigInputDto struct {
	Model                string
	ModelMaxTokens       int
	Temperature          float32
	TopP                 float32
	N                    int
	Stop                 []string
	MaxTokens            int
	PresencePenalty      float32
	FrequencyPenalty     float32
	InitialSystemMessage string
}

type ChatCompletionInputDto struct {
	ChatID      string
	UserID      string
	UserMessage string
	Config      ChatCompletionConfigInputDto
}

type ChatCompletionOutputDto struct {
	ChatID  string
	UserID  string
	Content string
}

type ChatCompletionUseCase struct {
	ChatGateway  gateway.ChatGateway
	OpenAiClient *openai.Client
	Stream       chan ChatCompletionOutputDto
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, openAIClient *openai.Client, stream chan ChatCompletionOutputDto) *ChatCompletionUseCase {
	return &ChatCompletionUseCase{
		ChatGateway:  chatGateway,
		OpenAiClient: openAIClient,
		Stream:       stream,
	}
}

func (usecase *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDto, stream chan ChatCompletionOutputDto) (*ChatCompletionOutputDto, error) {
	chat, err := usecase.ChatGateway.FindByChatID(ctx, input.ChatID)

	if err != nil {
		if err.Error() == "Chat not found" {
			chat, err := createNewChat(input)

			if err != nil {
				return nil, errors.New("Error creating chat: " + err.Error())
			}

			err = usecase.ChatGateway.CreateChat(ctx, chat)

			if err != nil {
				return nil, errors.New("Error persisting new chat: " + err.Error())
			}
		} else {
			return nil, errors.New("Error fetching existing chat: " + err.Error())
		}
	}

	userMessage, err := entity.NewMessage("user", input.UserMessage, chat.Config.Model)

	if err != nil {
		return nil, errors.New("Error creating user message: " + err.Error())
	}

	err = chat.AddMessage(userMessage)

	if err != nil {
		return nil, errors.New("Error adding new message: " + err.Error())
	}

	messages := []openai.ChatCompletionMessage{}

	for _, message := range chat.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	openaiResponse, err := usecase.OpenAiClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            chat.Config.Model.Name,
		Messages:         messages,
		MaxTokens:        chat.Config.MaxTokens,
		Temperature:      chat.Config.Temperature,
		TopP:             chat.Config.TopP,
		PresencePenalty:  chat.Config.PresencePenalty,
		FrequencyPenalty: chat.Config.FrequencyPenalty,
		Stop:             chat.Config.Stop,
		Stream:           true,
	})

	if err != nil {
		return nil, errors.New("Error creating chat completion: " + err.Error())
	}

	var fullResponse strings.Builder

	for {
		response, err := openaiResponse.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, errors.New("Error streaming response: " + err.Error())
		}

		fullResponse.WriteString(response.Choices[0].Delta.Content)

		r := ChatCompletionOutputDto{
			ChatID:  chat.ID,
			UserID:  input.UserID,
			Content: fullResponse.String(),
		}

		usecase.Stream <- r
	}

	assistant, err := entity.NewMessage("assistant", fullResponse.String(), chat.Config.Model)

	if err != nil {
		return nil, errors.New("Error creating assistant message: " + err.Error())
	}

	err = chat.AddMessage(assistant)

	if err != nil {
		return nil, errors.New("Error adding new message: " + err.Error())
	}

	err = usecase.ChatGateway.SaveChat(ctx, chat)

	if err != nil {
		return nil, errors.New("Error saving chat: " + err.Error())
	}

	return &ChatCompletionOutputDto{
		ChatID:  chat.ID,
		UserID:  input.UserID,
		Content: fullResponse.String(),
	}, nil
}

func createNewChat(input ChatCompletionInputDto) (*entity.Chat, error) {
	model := entity.NewModel(input.Config.Model, input.Config.ModelMaxTokens)
	chatConfig := &entity.ChatConfig{
		Temperature:      input.Config.Temperature,
		TopP:             input.Config.TopP,
		N:                input.Config.N,
		Stop:             input.Config.Stop,
		MaxTokens:        input.Config.MaxTokens,
		PresencePenalty:  input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
		Model:            model,
	}

	initialMessage, err := entity.NewMessage("system", input.Config.InitialSystemMessage, model)

	if err != nil {
		return nil, errors.New("Error creating initial message: " + err.Error())
	}

	chat, err := entity.NewChat(input.UserID, initialMessage, chatConfig)

	if err != nil {
		return nil, errors.New("Error creating new chat: " + err.Error())
	}

	return chat, nil
}
