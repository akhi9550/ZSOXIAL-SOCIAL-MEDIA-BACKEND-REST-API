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
	Caption   string    `json:"caption"`
	Tag       Tags      `json:"tag"`
	Url       string    `json:"url"`
	Likes     uint      `json:"likes"`
	Comments  uint      `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	Caption   string    `json:"caption"`
	Likes     uint      `json:"likes"`
	Comments  uint      `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type Url struct {
	ImageUrls []byte `json:"image_urls" gorm:"column:url"`
}

type Tags struct {
	User1 uint `json:"user1" gorm:"default:null"`
	User2 uint `json:"user2" gorm:"default:null"`
	User3 uint `json:"user3" gorm:"default:null"`
	User4 uint `json:"user4" gorm:"default:null"`
	User5 uint `json:"user5" gorm:"default:null"`
}

type Urls struct {
	ImageUrls []byte `json:"url"`
}

type UpdatePost struct {
	Caption string `json:"caption"`
	Tag     Tags   `json:"tag"`
}

type GetAllPosts struct {
	Post []PostResponse
}

type LikePostResponse struct {
	Id        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}

type PostCommentReq struct {
	PostID  uint   `json:"post_id"`
	Comment string `json:"comment"`
}

type PostCommentResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Profile   string    `json:"profile"`
	Comment   string    `json:"comment"`
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
