package models

import "time"

type NotificationReq struct {
	Offset string `query:"Offset" validate:"required"`
	Limit  string `query:"Limit" validate:"required"`
}

type Notification struct {
	UserID      int       `json:"user_id"`
	LikedUserID int       `json:"liked_user_id"`
	PostID      int       `json:"post_id"`
	Message     string    `json:"Message"`
	CreatedAt   time.Time `json:"created_at"`
}
