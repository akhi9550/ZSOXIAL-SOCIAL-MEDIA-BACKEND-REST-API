package di

import (
	server "github.com/akhi9550/auth-svc/pkg/api"
	"github.com/akhi9550/auth-svc/pkg/api/service"
	"github.com/akhi9550/auth-svc/pkg/client"
	"github.com/akhi9550/auth-svc/pkg/config"
	"github.com/akhi9550/auth-svc/pkg/db"
	"github.com/akhi9550/auth-svc/pkg/repository"
	"github.com/akhi9550/auth-svc/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	otpRepository := repository.NewOtpRepository(gormDB)
	adminRepository := repository.NewAdminRepository(gormDB)
	postClient := client.NewPostClient(&cfg)

	userUseCase := usecase.NewUserUseCase(userRepository)
	otpUseCase := usecase.NewOtpUseCase(otpRepository)
	adminUseCase := usecase.NewAdminUseCase(adminRepository,postClient)

	ServiceServer := service.NewAuthServer(userUseCase, adminUseCase, otpUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, ServiceServer)
	if err != nil {
		return &server.Server{}, err
	}
	return grpcServer, nil
}
