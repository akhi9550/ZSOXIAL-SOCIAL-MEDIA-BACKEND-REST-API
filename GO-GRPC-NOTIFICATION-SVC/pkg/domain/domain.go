package domain

import "time"

type NotificationType string

const (
	CommentNotification NotificationType = "comment"
	FollowNotification  NotificationType = "follow"
)

type LikeNotification struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	RecipientID      uint             `json:"recipient_id"`
	SenderID         uint             `json:"sender_id"`
	NotificationType NotificationType `json:"notification_type"`
	PostID           uint             `json:"post_id"`
	CommentID        uint             `json:"comment_id"`
	Seen             bool             `json:"seen" gorm:"default:false"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type Notification struct {
	ID      int64  `json:"id" gorm:"primaryKey"`
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
	PostID  int64  `json:"post_id"`
}
