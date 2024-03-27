package di

import (
	server "github.com/akhi9550/pkg/api"
	"github.com/akhi9550/pkg/api/service"
	"github.com/akhi9550/pkg/config"
	"github.com/akhi9550/pkg/db"
	"github.com/akhi9550/pkg/repository"
	"github.com/akhi9550/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	otpRepository := repository.NewOtpRepository(gormDB)
	adminRepository := repository.NewAdminRepository(gormDB)

	userUseCase := usecase.NewUserUseCase(userRepository)
	otpUseCase := usecase.NewOtpUseCase(otpRepository)
	adminUseCase := usecase.NewAdminUseCase(adminRepository)

	ServiceServer := service.NewAuthServer(userUseCase, adminUseCase, otpUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, ServiceServer)
	if err != nil {
		return &server.Server{}, err
	}
	return grpcServer, nil
}
