package domain

import "time"

type Story struct {
	ID         uint      `json:"id" gorm:"uniquekey; not null"`
	UserID     uint      `json:"user_id"`
	Url        string    `json:"url"`
	LikesCount uint      `json:"likes_count" gorm:"default:0"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	IsValid    bool      `json:"is_valid" gorm:"default:false"`
}

type StoryLike struct {
	ID        uint      `json:"id" gorm:"uniquekey; not null"`
	StoryID   uint      `json:"story_id"`
	Story     Story     `json:"story" gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE"`
	LikedUser uint      `json:"liked_user"`
	CreatedAt time.Time `json:"created_at"`
}

type ViewStory struct {
	ID       uint
	StoryID  uint   `json:"story_id"`
	Story    Story  `json:"story" gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE"`
	ViewerID uint   `json:"viewer_id"`
	Viewer   string `json:"viewer"`
}

type ArchiveStory struct {
	ID      uint `json:"id" gorm:"uniquekey; not null"`
	UserID  uint `json:"user_id"`
	StoryID uint `json:"story_id"`
}
