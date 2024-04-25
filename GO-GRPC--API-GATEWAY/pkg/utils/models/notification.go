package models

import "time"

type NotificationRequest struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}

type NotificationReq struct {
	Offset string `query:"Offset" validate:"required"`
	Limit  string `query:"Limit" validate:"required"`
}

type Notification struct {
	LikedUser   string    `json:"liked_user"`
	UserProfile string    `json:"user_profile"`
	PostID      int       `json:"post_id"`
	Content     string    `json:"Content"`
	Timestamp   time.Time `json:"TimeStamp"`
}

type Response struct {
	Message string `json:"message"`
}

type Responses struct {
	UserID  uint   `json:"user_id"`
	PostID  uint   `json:"post_id"`
	Message string `json:"message"`
	Content string `json:"content"`
}

type ConsumeResponses struct {
	UserID  uint   `json:"user_id"`
	PostID  uint   `json:"post_id"`
	Message string `json:"message"`
}
