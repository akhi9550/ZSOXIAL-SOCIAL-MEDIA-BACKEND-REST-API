package models

type NotificationRequest struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
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
