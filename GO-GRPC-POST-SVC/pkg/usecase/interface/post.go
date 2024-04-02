package interfaces

import "github.com/akhi9550/post-svc/pkg/utils/models"

type PostUseCase interface {
	CreatePost(userID int, data models.PostRequest, file []byte, users []models.Tag) (models.PostResponse, error)
	GetPost(userID int, postID int) (models.PostResponse, error)
	UpdatePost(userID int, data models.UpdatePostReq,tag []models.Tag) (models.UpdateResponse, error)
	DeletePost(userID int, postID int) error
	GetAllPost(userID int) ([]models.PostResponse, error)
	// LikePost(userID int, postID int) (models.LikePostResponse, error)
	// PostComment(userID int, data models.PostCommentReq) (models.PostCommentResponse, error)
	// GetComment(postID int) (models.GetCommentResponse, error)
	// ReplyComment(userID int, data models.CommentReply) error
	// DeleteComment(userID int, commentID int) error
}
