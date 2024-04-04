package domain

import "time"

type Post struct {
	ID            uint      `json:"id" gorm:"uniquekey; not null"`
	UserID        uint      `json:"user_id"`
	Url           string    `json:"url"`
	Caption       string    `json:"caption"`
	TypeID        uint      `json:"type_id"`
	LikesCount    uint      `json:"likes_count" gorm:"default:0"`
	CommentsCount uint      `json:"comments_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	IsArchive     bool      `json:"is_archive" gorm:"default:false"`
}

type Tags struct {
	UserID  uint `json:"user_id"`
	PostID  uint `json:"post_id"`
	Taguser uint `json:"taguser" gorm:"default:null"`
}

type PostImages struct {
	ID     uint   `json:"id" gorm:"uniquekey; not null"`
	PostID uint   `json:"post_id"`
	Urls   string `json:"url"`
}

type PostType struct {
	ID   uint   `json:"id" gorm:"uniquekey; not null"`
	Type string `json:"type"`
	// gorm:"type:post;default:'post';check:type IN ('reel', 'post', 'video')"
}

type Likes struct {
	ID        uint      `json:"id" gorm:"uniquekey; not null"`
	PostID    uint      `json:"post_id"`
	Post      Post      `json:"post" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	LikedUser uint      `json:"liked_user"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID            uint      `json:"id" gorm:"uniquekey; not null"`
	PostID        uint      `json:"post_id"`
	Post          Post      `json:"post" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentedUser uint      `json:"commented_user"`
	CommentData   string    `json:"comment_data"`
	CreatedAt     time.Time `json:"created_at"`
}

type CommentReplies struct {
	ID            uint      `json:"id" gorm:"uniquekey; not null"`
	Post_id       uint      `json:"post_id"`
	CommentID     uint      `json:"comment_id"`
	CommentedUser uint      `json:"commented_user"`
	ReplyUser     uint      `json:"reply_user"`
	Replies       string    `json:"replies"`
	CreatedAt     time.Time `json:"created_at"`
}

type SavedPost struct {
	ID     uint `json:"id" gorm:"uniquekey; not null"`
	PostID uint `json:"post_id"`
	UserId uint `json:"user_id"`
}

type ArchivePost struct {
	ID     uint `json:"id" gorm:"uniquekey; not null"`
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}
