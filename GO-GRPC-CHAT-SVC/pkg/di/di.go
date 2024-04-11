package di

import (
	server "github.com/akhi9550/chat-svc/pkg/api"
	"github.com/akhi9550/chat-svc/pkg/api/service"
	"github.com/akhi9550/chat-svc/pkg/client"
	"github.com/akhi9550/chat-svc/pkg/config"
	"github.com/akhi9550/chat-svc/pkg/db"
	"github.com/akhi9550/chat-svc/pkg/repository"
	"github.com/akhi9550/chat-svc/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {
	database, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	chatRepository := repository.NewChatRepository(database)
	authClient := client.NewAuthClient(&cfg)

	chatUseCase := usecase.NewChatUseCase(chatRepository, authClient)

	ServiceServer := service.NewPostServer(chatUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, ServiceServer)
	if err != nil {
		return &server.Server{}, err
	}
	return grpcServer, nil
}
