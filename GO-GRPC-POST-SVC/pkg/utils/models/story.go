package models

import (
	"time"
)

// type CreateStory struct {
// 	UserID   uint   `json:"user_id"`
// 	ImageUrl string `json:"image_url"`
// }
// type StoryResponse struct {
// 	Author    UserData  `json:"author"`
// 	ImageUrl  string    `json:"image_url"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type StoryUrl struct {
// 	Imageurl []byte `json:"imageurl" gorm:"validate:required"`
// }

type CreateStoryResponse struct {
	Author    UserData  `json:"author"`
	Story     string    `json:"story"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStory struct {
	Story     string    `json:"story"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStories struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	StartTime time.Time `json:"start_time"`
}

type CreateStoryResponses struct {
	Author    UserData  `json:"author"`
	StoryID   uint      `json:"story_id"`
	Story     string    `json:"story"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStoriesResponse struct {
	StoryID   uint      `json:"story_id" gorm:"column:id"`
	Story     string    `json:"story" gorm:"column:url"`
	CreatedAt time.Time `json:"created_at" gorm:"column:start_time"`
}

type Likeduser struct {
	LikeUserID uint `json:"like_user_id" gorm:"column:liked_user"`
}

type Likedusers struct {
	LikeUser string `json:"like_user"`
	Profile  string `json:"profile"`
}

type Viewer struct {
	ViewerID uint `json:"viewer_id"`
}

type Viewers struct {
	ViewUser string `json:"view_user"`
	Profile  string `json:"profile"`
}

type StoryDetails struct {
	StoryID   uint         `json:"story_id"`
	LikedUser []Likedusers `json:"liked_user"`
	Viewer    []Viewers    `json:"viewer"`
}
