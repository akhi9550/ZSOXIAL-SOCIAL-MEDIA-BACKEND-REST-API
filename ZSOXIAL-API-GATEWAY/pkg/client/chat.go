package client

import (
	"fmt"

	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/chat"
	"google.golang.org/grpc"
)

type ChatClient struct {
	Client pb.ChatServiceClient
}

func NewChatClient(cfg config.Config) *ChatClient {
	grpcConnection, err := grpc.Dial(cfg.ChatSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewChatServiceClient(grpcConnection)

	return &ChatClient{
		Client: grpcClient,
	}
}
