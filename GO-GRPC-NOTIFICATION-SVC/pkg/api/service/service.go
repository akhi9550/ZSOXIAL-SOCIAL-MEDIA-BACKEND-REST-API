package service

import (
	"context"

	interfaces"github.com/akhi9550/notification-svc/pkg/usecase/interface"
	pb"github.com/akhi9550/notification-svc/pkg/pb/notification"
	"github.com/akhi9550/notification-svc/pkg/utils/models"
)

type NotificationHandler struct {
	notificationUsecase interfaces.NotificationUsecaseInterface
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationServer(UseCase interfaces.NotificationUsecaseInterface) pb.NotificationServiceServer {
	return &NotificationHandler{
		notificationUsecase: UseCase,
	}
}
func (n *NotificationHandler) SendLikeNotification(ctx context.Context, notify *pb.LikeNotification) (*pb.NotificationResponse, error) {
	_, err := n.notificationUsecase.AddLikeNotification(models.LikeNotification{
		UserID: notify.UserId,
		PostID: notify.PostId,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.NotificationResponse{
		Message: "successfully got all notifications",
	}

	return response, nil

}

func (n *NotificationHandler) ConsumeKafkaCommentMessages(ctx context.Context, p *pb.ConsumeKafkaCommentMessagesRequest) (*pb.ConsumeKafkaCommentMessagesResponse, error) {

	res, err := n.notificationUsecase.ConsumeCommentMessage(p.UserId)

	if err != nil {
		return &pb.ConsumeKafkaCommentMessagesResponse{}, err
	}

	return &pb.ConsumeKafkaCommentMessagesResponse{
		UserId:  res.UserID,
		PostId:  res.PostID,
		Message: res.Message,
		Content: res.Content,
	}, nil

}

func (n *NotificationHandler) ConsumeKafkaLikeMessages(ctx context.Context, p *pb.ConsumeKafkaLikeMessagesRequest) (*pb.ConsumeKafkaLikeMessagesResponse, error) {

	response, err := n.notificationUsecase.ConsumeMessage(p.UserId)

	if err != nil {
		return &pb.ConsumeKafkaLikeMessagesResponse{}, err
	}
	return &pb.ConsumeKafkaLikeMessagesResponse{
		UserId:  response.UserID,
		PostId:  response.PostID,
		Message: response.Message,
		Content: response.Content,
	}, nil

}
