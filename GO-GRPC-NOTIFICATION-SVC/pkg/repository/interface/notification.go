package interfaces

import (
	"github.com/akhi9550/notification-svc/pkg/utils/models"
)

type NotificationRepository interface {
	StoreNotification(models.NotificationReq) error
	GetNotification(userID int, req models.Pagination) ([]models.Notification, error)
}
