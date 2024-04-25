package domain

import "time"

type Notification struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id"`
	LikedUserID int       `json:"liked_user_id"`
	PostID      int       `json:"post_id"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}
