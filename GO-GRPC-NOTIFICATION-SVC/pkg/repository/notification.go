package repository

import (
	interfaces "github.com/akhi9550/notification-svc/pkg/repository/interface"
	"github.com/akhi9550/notification-svc/pkg/utils/models"
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
func (n *notificationRepository) StoreNotification(notification models.NotificationReq) error {
	err := n.DB.Exec(`INSERT INTO notifications(user_id,sender_id,post_id,message,created_at) VALUES (?,?,?,?,?)`, notification.UserID, notification.SenderID, notification.PostID, notification.Message, notification.CreatedAt).Error
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationRepository) GetNotification(userID int, pagination models.Pagination) ([]models.Notification, error) {
	var data []models.Notification
	if pagination.Offset <= 0 {
		pagination.Offset = 1
	}
	offset := (pagination.Offset - 1) * pagination.Limit
	err := n.DB.Raw(`SELECT sender_id, message, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, userID, pagination.Limit, offset).Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
