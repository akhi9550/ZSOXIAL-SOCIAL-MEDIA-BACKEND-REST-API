package repository

import (
	"github.com/akhi9550/notification-svc/pkg/domain"
	interfaces"github.com/akhi9550/notification-svc/pkg/repository/interface"
	"gorm.io/gorm"
)

type notificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) interfaces.NotificationRepository {
	return &notificationRepository{
		DB: db,
	}
}
func (n *notificationRepository) AddNotification(notification domain.Notification) (int64, error) {
	err :=n.DB.Exec(`INSERT INTO notifications (user_id, message, post_id) VALUES (?, ?, ?)`,notification.UserID,notification.Message,notification.PostID).Error
	if err != nil {
		return 0, err
	}
	return notification.UserID, nil
}
