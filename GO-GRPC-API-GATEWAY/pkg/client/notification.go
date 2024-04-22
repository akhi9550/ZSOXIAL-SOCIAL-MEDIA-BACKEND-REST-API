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

func (n *NotificationClient) SendLikeNotification(userID, postID int) (string, error) {
	data, err := n.Client.SendLikeNotification(context.Background(), &pb.LikeNotification{
		PostId: int64(postID),
		UserId: int64(userID),
	})
	if err != nil {
		return "", err
	}
	return data.Message, nil
}

func (n *NotificationClient) ConsumeKafkaLikeMessages(userID int) (models.Responses, error) {
	data, err := n.Client.ConsumeKafkaLikeMessages(context.Background(), &pb.ConsumeKafkaLikeMessagesRequest{
		UserId: int64(userID),
	})
	if err != nil {
		return models.Responses{}, err
	}
	return models.Responses{
		UserID:  uint(data.UserId),
		PostID:  uint(data.PostId),
		Message: data.Message,
		Content: data.Content,
	}, nil
}

func (n *NotificationClient) ConsumeKafkaCommentMessages(userID int) (models.Responses, error) {
	data, err := n.Client.ConsumeKafkaCommentMessages(context.Background(), &pb.ConsumeKafkaCommentMessagesRequest{
		UserId: int64(userID),
	})
	if err != nil {
		return models.Responses{}, err
	}
	return models.Responses{
		UserID:  uint(data.UserId),
		PostID:  uint(data.PostId),
		Message: data.Message,
		Content: data.Content,
	}, nil
}
 