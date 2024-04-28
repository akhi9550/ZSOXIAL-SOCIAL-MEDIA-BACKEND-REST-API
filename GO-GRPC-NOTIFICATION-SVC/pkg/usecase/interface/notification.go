package interfaces

import "github.com/akhi9550/notification-svc/pkg/utils/models"

type NotificationUsecaseInterface interface {
	ConsumeLikeMessage()
	ConsumeCommentMessage()
	ConsumeFollowReqMessage()
	ConsumeAcceptFollowReqMessage()
	GetNotification(int, models.Pagination) ([]models.NotificationResponse, error)
}
