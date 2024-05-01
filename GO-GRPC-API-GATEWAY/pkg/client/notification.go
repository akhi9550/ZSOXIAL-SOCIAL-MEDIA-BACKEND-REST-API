package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/notification"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
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

func (n *NotificationClient) GetNotification(userID int, req models.NotificationPagination) ([]models.NotificationResponse, error) {
data,err:=n.Client.GetNotification(context.Background(),&pb.GetNotificationRequest{
	UserID: int64(userID),
	Limit: int64(req.Limit),
	Offset: int64(req.Offset),
})
if err!=nil{
	return []models.NotificationResponse{},err
}
var response []models.NotificationResponse
	for _, v := range data.Notification {
		notificationResponse := models.NotificationResponse{
			UserID:    int(v.UserID),
			Username:  v.Username,
			Profile:   v.Profile,
			Message:   v.Message,
			CreatedAt: v.Time,
		}
		response = append(response, notificationResponse)
	}
	return response, nil
}
