package di

import (
	"sync"

	server "github.com/akhi9550/notification-svc/pkg/api"
	"github.com/akhi9550/notification-svc/pkg/api/service"
	"github.com/akhi9550/notification-svc/pkg/client"
	"github.com/akhi9550/notification-svc/pkg/config"
	"github.com/akhi9550/notification-svc/pkg/db"
	"github.com/akhi9550/notification-svc/pkg/repository"
	"github.com/akhi9550/notification-svc/pkg/usecase"
)

var configMutex sync.Mutex

func InitializeAPI(cfg config.Config) (*server.Server, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	notificationRepository := repository.NewNotificationRepository(gormDB)
	authClient := client.NewAuthClient(&cfg)
	notificationUseCase := usecase.NewUserUserUseCase(notificationRepository, authClient)
	ServiceServer := service.NewNotificationServer(notificationUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, ServiceServer)
	if err != nil {
		return &server.Server{}, err
	}
	go notificationUseCase.ConsumeNotification()
	return grpcServer, nil
}
