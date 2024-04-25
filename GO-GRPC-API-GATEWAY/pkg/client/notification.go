package client

import (
	"fmt"

	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/notification"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	Client pb.NotificationServiceClient
}

func NewNotificationClient(cfg config.Config) *NotificationClient {
	grpcConnection, err := grpc.Dial(cfg.NotificationSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewNotificationServiceClient(grpcConnection)

	return &NotificationClient{
		Client: grpcClient,
	}
}
