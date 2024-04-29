package interfaces

import "github.com/akhi9550/api-gateway/pkg/utils/models"

type NotificationClient interface {
	GetNotification(userID int, req models.NotificationPagination) ([]models.NotificationResponse, error)
}
