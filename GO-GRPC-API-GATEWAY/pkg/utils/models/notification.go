package models

import "time"

type NotificationRequest struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}

type NotificationPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Notification struct {
	UserID      int       `json:"user_id"`
	LikedUserID int       `json:"liked_user_id"`
	PostID      int       `json:"post_id"`
	Message     string    `json:"Message"`
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
