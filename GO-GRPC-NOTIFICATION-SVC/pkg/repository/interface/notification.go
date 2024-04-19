package interfaces

import "github.com/akhi9550/notification-svc/pkg/domain"

type NotificationRepository interface {
	AddNotification(notification domain.Notification) (int64, error)
}
