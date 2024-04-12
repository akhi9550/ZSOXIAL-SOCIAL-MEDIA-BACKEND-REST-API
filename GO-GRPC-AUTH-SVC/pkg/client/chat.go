package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/auth-svc/pkg/config"
	pb "github.com/akhi9550/auth-svc/pkg/pb/chat"
	"google.golang.org/grpc"
)

type clientChat struct {
	Client pb.ChatServiceClient
}

func NewChatClient(c *config.Config) *clientChat {
	cc, err := grpc.Dial(c.ChatSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	pbClient := pb.NewChatServiceClient(cc)
	
	return &clientChat{
		Client: pbClient,
	}
}

func (ch *clientChat) CreateChatRoom(user1, user2 int) error {
	_, err := ch.Client.CreateChatRoom(context.Background(), &pb.CreateChatRoomRequest{
		Userid:      int64(user1),
		Followingid: int64(user2),
	})
	if err != nil {
		fmt.Println("data", err)
		return err
	}
	return nil
}
