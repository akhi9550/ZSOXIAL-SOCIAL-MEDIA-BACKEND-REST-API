package di

import (
	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/client"
	"github.com/akhi9550/api-gateway/pkg/config"
	server "github.com/akhi9550/api-gateway/pkg/api"
)

func InitializeAPI(cfg config.Config) (*server.ServerHTTP, error) {
	authClient := client.NewAuthClient(cfg)
	authHandler := handler.NewAuthHandler(authClient)

	serverHTTP := server.NewServerHTTP(authHandler)
	return serverHTTP, nil
}
