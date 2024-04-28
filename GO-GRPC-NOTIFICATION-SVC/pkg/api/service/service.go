package service

import (
	"context"

	pb "github.com/akhi9550/notification-svc/pkg/pb/notification"
	interfaces "github.com/akhi9550/notification-svc/pkg/usecase/interface"
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

func (n *NotificationHandler) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	result, err := n.notificationUsecase.GetNotification(int(req.UserID), models.Pagination{Limit: int(req.Limit), Offset: int(req.Offset)})
	if err != nil {
		return nil, err
	}

	var finalResult []*pb.Message
	for _, val := range result {
		finalResult = append(finalResult, &pb.Message{
			UserID:   int64(val.UserID),
			Username: val.Username,
			Profile:  val.Profile,
			PostID:   int64(val.PostID),
			Message:  val.Message,
			Time:     val.CreatedAt,
		})
	}
	return &pb.GetNotificationResponse{Notification: finalResult}, nil
}
