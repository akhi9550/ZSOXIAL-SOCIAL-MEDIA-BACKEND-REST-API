package models

import "time"

type Notification struct {
	UserID    int       `json:"user_id"`
	SenderID  int       `json:"sender_id"`
	PostID    int       `json:"post_id"`
	Message   string    `json:"Message"`
	CreatedAt time.Time `json:"created_at"`
}
