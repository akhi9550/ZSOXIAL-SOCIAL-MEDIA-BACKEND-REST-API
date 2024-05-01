package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/chat"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
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

func (c *ChatClient) GetChat(userID string, req models.ChatRequest) ([]models.TempMessage, error) {
	data, err := c.Client.GetFriendChat(context.Background(), &pb.GetFriendChatRequest{
		UserID:   userID,
		FriendID: req.FriendID,
		OffSet:   req.Offset,
		Limit:    req.Limit,
	})
	if err != nil {
		return []models.TempMessage{}, err
	}
	var response []models.TempMessage
	for _, v := range data.FriendChat {
		chatResponse := models.TempMessage{
			SenderID:    v.SenderId,
			RecipientID: v.RecipientId,
			Content:     v.Content,
			Timestamp:   v.Timestamp,
		}
		response = append(response, chatResponse)

	}
	return response, nil
}
