package domain

import "time"

type Story struct {
	ID        uint      `json:"id" gorm:"uniquekey; not null"`
	UserID    uint      `json:"user_id"`
	Data      string    `json:"data"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsValid   bool      `json:"is_valid"`
}

type StoryLike struct {
	ID        uint      `json:"id" gorm:"uniquekey; not null"`
	StoryID   uint      `json:"story_id"`
	Story     Story     `json:"story" gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE"`
	LikedUser uint      `json:"liked_user"`
	CreatedAt time.Time `json:"created_at"`
}

type ViewStory struct {
	ID      uint
	StoryID uint  `json:"story_id"`
	Story   Story `json:"story" gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE"`
	Viewer  uint  `json:"viewer"`
}
