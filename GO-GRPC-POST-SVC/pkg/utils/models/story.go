package models

import (
	"time"
)

type CreateStory struct {
	UserID   uint   `json:"user_id"`
	ImageUrl string `json:"image_url"`
}
type StoryResponse struct {
	Author    UserData  `json:"author"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type StoryUrl struct {
	Imageurl []byte `json:"imageurl" gorm:"validate:required"`
}
