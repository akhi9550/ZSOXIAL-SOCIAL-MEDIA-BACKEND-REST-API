package di

import (
	"fmt"

	server "github.com/akhi9550/api-gateway/pkg/api"
	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/client"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/helper"
	pb "github.com/akhi9550/api-gateway/pkg/pb/auth"
	"github.com/redis/go-redis/v9"
)

func InitializeAPI(cfg config.Config) (*server.ServerHTTP, error) {
	authClient := client.NewAuthClient(cfg)
	authCachig := AuthCache(cfg)
	authHandler := handler.NewAuthHandler(authClient, authCachig)

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

var authclient pb.AuthServiceClient

func AuthCache(cfg config.Config) *helper.RedisCaching {
	return helper.NewRedisCaching(InitRedisDB(cfg), authclient)
}

func InitRedisDB(config config.Config) *redis.Client {
	fmt.Println("cccccccccc", config.RedisUrl, config.RedisPassword)
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		// DB:       0,
		// Username: "default",
	})
	return client
}
