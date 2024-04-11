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
	Likes     uint      `json:"likes"`
	Comments  uint      `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

type Tags struct {
	User1 uint `json:"user1"`
	User2 uint `json:"user2"`
	User3 uint `json:"user3"`
	User4 uint `json:"user4"`
	User5 uint `json:"user5"`
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
	User string `json:"user"`
}
type Url struct {
	ImageUrls string `json:"image_urls"`
}

type UpdatePost struct {
	Caption string `json:"caption"`
	Tag     Tags   `json:"tag"`
}

type GetAllPosts struct {
	Post []PostResponse
}

type LikePostResponse struct {
	UserID    uint      `json:"user_id"`
	LikedUser string    `json:"like_user"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}

type PostCommentReq struct {
	PostID  uint   `json:"post_id"`
	Comment string `json:"comment"`
}

type ReplyCommentReq struct {
	CommentID uint   `json:"comment_id"`
	Reply     string `json:"reply"`
}

type PostCommentResponse struct {
	UserID      uint      `json:"user_id"`
	CommentUser string    `json:"comment_user"`
	Profile     string    `json:"profile"`
	CommentID   uint      `json:"comment_id"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostComment struct {
	UserID      uint      `json:"user_id"`
	CommentUser string    `json:"commented_user"`
	Profile     string    `json:"profile"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReplyPostCommentResponse struct {
	UserID    uint      `json:"user_id"`
	ReplyUser string    `json:"reply_user"`
	Profile   string    `json:"profile"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `json:"created_at"`
}

type ReplyReposne struct {
	Comment PostComment
	Reply   ReplyPostCommentResponse
}

type AllComments struct {
	CommentUser string    `json:"commented_user"`
	Profile     string    `json:"profile"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type AllReplies struct {
	ReplyUser string    `json:"reply_user"`
	Profile   string    `json:"profile"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `json:"created_at"`
}

type AllCommentsAndReplies struct {
	CommentUser string    `json:"commented_user"`
	Profile     string    `json:"profile"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
	Reply       []AllReplies
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

type CreateStoryResponse struct {
	Author    UserData  `json:"author"`
	Story     string    `json:"story"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStoryResponses struct {
	Author    UserData  `json:"author"`
	StoryID   uint      `json:"story_id"`
	Story     string    `json:"story"`
	CreatedAt time.Time `json:"created_at"`
}

type ReportPostRequest struct {
	PostID uint   `json:"post_id"`
	Report string `json:"report"`
}

type Likeduser struct {
	LikeUser string `json:"like_user"`
	Profile  string `json:"profile"`
}

type Viewer struct {
	ViewUser string `json:"view_user"`
	Profile  string `json:"profile"`
}

type StoryDetails struct {
	StoryID   uint        `json:"story_id"`
	LikedUser []Likeduser `json:"liked_user"`
	Viewer    []Viewer    `json:"viewer"`
}
