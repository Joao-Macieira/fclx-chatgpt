package service

import (
	"github.com/joao-macieira/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/joao-macieira/fclx/chatservice/internal/usecases/chatcompletionstream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigInputDto
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDto
}

func NewChatService(chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDto, streamChannel chan chatcompletionstream.ChatCompletionOutputDto) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}
}

func (chat *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompletionstream.ChatCompletionConfigInputDto{
		Model:                chat.ChatConfigStream.Model,
		ModelMaxTokens:       chat.ChatConfigStream.ModelMaxTokens,
		Temperature:          chat.ChatConfigStream.Temperature,
		TopP:                 chat.ChatConfigStream.TopP,
		N:                    chat.ChatConfigStream.N,
		Stop:                 chat.ChatConfigStream.Stop,
		MaxTokens:            chat.ChatConfigStream.MaxTokens,
		InitialSystemMessage: chat.ChatConfigStream.InitialSystemMessage,
	}

	input := chatcompletionstream.ChatCompletionInputDto{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()

	go func() {
		for msg := range chat.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}
	}()

	_, err := chat.ChatCompletionStreamUseCase.Execute(ctx, input)

	if err != nil {
		return err
	}

	return nil
}
