package interfaces

import "github.com/akhi9550/notification-svc/pkg/utils/models"

type NotificationUsecaseInterface interface {
	AddLikeNotification(notification models.LikeNotification) (int64, error)
	ConsumeMessage(user_id int64) (models.LikeNotification, error)
	ConsumeCommentMessage(user_id int64) (models.CommentNotification, error)
}
