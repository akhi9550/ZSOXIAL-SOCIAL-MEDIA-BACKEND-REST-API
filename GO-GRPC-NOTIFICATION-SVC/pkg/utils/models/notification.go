
package models

type LikeNotification struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
	PostID  int64  `json:"post_id"`
	Content string `json:"content"`
}
type CommentNotification struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
	PostID  int64  `json:"post_id"`
	Content string `json:"content"`
}
