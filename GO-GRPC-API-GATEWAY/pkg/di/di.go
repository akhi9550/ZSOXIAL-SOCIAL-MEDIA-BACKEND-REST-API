package di

import (
	server "github.com/akhi9550/api-gateway/pkg/api"
	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/client"
	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	pbnotification "github.com/akhi9550/api-gateway/pkg/pb/notification"
	pbchat "github.com/akhi9550/api-gateway/pkg/pb/chat"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/redis/go-redis/v9"
)

func InitializeAPI(cfg config.Config) (*server.ServerHTTP, error) {
	authClient := client.NewAuthClient(cfg)
	authCachig := AuthCache(cfg, authClient)
	authHandler := handler.NewAuthHandler(authClient, authCachig)

	postClient := client.NewPostClient(cfg)
	postCachig := PostCache(cfg, postClient)
	postHandler := handler.NewPostHandler(postClient, postCachig)

	helper := helper.NewHelper(&cfg)
	chatClient := client.NewChatClient(cfg)
	chatCaching:=ChatCache(cfg)
	chatHandler := handler.NewChatHandler(chatClient.Client, helper,chatCaching)

	notificationClient := client.NewNotificationClient(cfg)
	notificationCaching:= NotificationCache(cfg)
	notificationHandler := handler.NewNotificationHandler(notificationClient.Client,notificationCaching)

	serverHTTP := server.NewServerHTTP(authHandler, postHandler, chatHandler, notificationHandler)
	return serverHTTP, nil
}

func AuthCache(cfg config.Config, authClient interfaces.AuthClient) *helper.RedisAuthCaching {
	return helper.NewRedisAuthCaching(InitRedisDB(cfg), authClient)
}

func PostCache(cfg config.Config, postClient interfaces.PostClient) *helper.RedisPostCaching {
	return helper.NewRedisPostCaching(InitRedisDB(cfg), postClient)
}

var notificationClient pbnotification.NotificationServiceClient
func NotificationCache(cfg config.Config) *helper.RedisNotificationCaching {
	return helper.NewRedisNotificationCaching(InitRedisDB(cfg), notificationClient)
}

var chatClient pbchat.ChatServiceClient
func ChatCache(cfg config.Config) *helper.RedisChatCaching {
	return helper.NewRedisChatCaching(InitRedisDB(cfg), chatClient)
}

func InitRedisDB(config config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		// DB:       0,
		// Username: "default",
	})
	return client
}
