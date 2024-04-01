package di

import (
	server "github.com/akhi9550/post-svc/pkg/api"
	"github.com/akhi9550/post-svc/pkg/api/service"
	"github.com/akhi9550/post-svc/pkg/client"
	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/akhi9550/post-svc/pkg/db"
	"github.com/akhi9550/post-svc/pkg/repository"
	"github.com/akhi9550/post-svc/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	postRepository := repository.NewPostRepository(gormDB)
	authClient := client.NewAuthClient(&cfg)
	// storyRepository := repository.NewStoryRepository(gormDB)

	postUseCase := usecase.NewPostUseCase(postRepository, authClient)
	// storyUseCase := usecase.NewStoryUseCase(storyRepository)

	ServiceServer := service.NewPostServer(postUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, ServiceServer)
	if err != nil {
		return &server.Server{}, err
	}
	return grpcServer, nil
}
