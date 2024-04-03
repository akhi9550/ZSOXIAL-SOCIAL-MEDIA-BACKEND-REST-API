package models

import (
	"time"
)

type PostRequest struct {
	Caption string `json:"caption"`
	TypeId  uint   `json:"TypeId"`
}
type UserData struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Profile  string `json:"profile"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	Author    UserData  `json:"author"`
	Tag       []Tag     `json:"tag"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes"`
	Comments  uint      `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type ArchivePostResponse struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes" gorm:"column:likes_count"`
	Comments  uint      `json:"comments" gorm:"column:comments_count"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes"  gorm:"column:likes_count"`
	Comments  uint      `json:"comments"  gorm:"column:comments_count"`
	CreatedAt time.Time `json:"created_at"`
}

type SavedResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes"  gorm:"column:likes_count"`
	Comments  uint      `json:"comments"  gorm:"column:comments_count"`
	CreatedAt time.Time `json:"created_at"`
}

type Url struct {
	ImageUrls []byte `json:"image_urls" gorm:"column:url"`
}

type Urls struct {
	ImageUrls []byte `json:"url"`
}

type GetAllPosts struct {
	Post []PostResponse
}

type LikePostResponse struct {
	UserID    uint      `json:"user_id"`
	LikedUser string    `json:"liked_user"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}

type LikesReponse struct {
	UserID    uint      `json:"user_id" gorm:"column:liked_user"`
	CreatedAt time.Time `json:"created_at"`
}

type PostCommentReq struct {
	PostID  uint   `json:"post_id"`
	Comment string `json:"comment"`
}

type PostCommentResponse struct {
	UserID        uint      `json:"user_id"`
	CommentedUser string    `json:"commented_user"`
	Profile       string    `json:"profile"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostCommentResponses struct {
	UserID    uint      `json:"user_id" gorm:"column:commented_user"`
	Comment   string    `json:"comment" gorm:"column:comment_data"`
	CreatedAt time.Time `json:"created_at"`
}

type Comments struct {
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
type GetCommentResponse struct {
	PostID  uint       `json:"post_id"`
	Comment []Comments `json:"comment"`
}

type CommentReply struct {
	CommentID uint   `json:"comment_id"`
	PostID    uint   `json:"post_id"`
	Reply     string `json:"comment"`
}

type UpdatePostReq struct {
	PostID  uint   `json:"post_id"`
	Caption string `json:"caption"`
	TypeID  uint   `json:"type_id"`
}

type UpdateResponse struct {
	ID        uint      `json:"id"`
	Author    UserData  `json:"author"`
	Tag       []Tag     `json:"tag"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes"`
	Comments  uint      `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	User string `json:"user" gorm:"column:taguser"`
}

type PostID struct {
	PostID int `json:"post_id"`
}

type TagUsers struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}
