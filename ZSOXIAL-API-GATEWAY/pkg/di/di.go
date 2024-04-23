package di

import (
	server "github.com/akhi9550/api-gateway/pkg/api"
	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/client"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/helper"
)

func InitializeAPI(cfg config.Config) (*server.ServerHTTP, error) {
	authClient := client.NewAuthClient(cfg)
	authHandler := handler.NewAuthHandler(authClient)

	postClient := client.NewPostClient(cfg)
	postHandler := handler.NewPostHandler(postClient)

	helper := helper.NewHelper(&cfg)

	chatClient := client.NewChatClient(cfg)
	chatHandler := handler.NewChatHandler(chatClient.Client, helper)

	notificationClient := client.NewNotificationClient(cfg)
	notificationHandler := handler.NewNotificationHandler(notificationClient.Client)
	
	serverHTTP := server.NewServerHTTP(authHandler, postHandler, chatHandler, notificationHandler)
	return serverHTTP, nil
}
