package server

import (
	"net"

	"github.com/joao-macieira/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/joao-macieira/fclx/chatservice/internal/infra/grpc/service"
	"github.com/joao-macieira/fclx/chatservice/internal/usecases/chatcompletionstream"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigInputDto
	ChatService                 service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDto
}

func NewGRPCServer(chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDto, port, authToken string, streamChannel chan chatcompletionstream.ChatCompletionOutputDto) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)
	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               streamChannel,
		ChatService:                 *chatService,
	}
}

func (g *GRPCServer) Start() {
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &g.ChatService)

	lis, err := net.Listen("tcp", ":"+g.Port)

	if err != nil {
		panic(err.Error())
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}
}
