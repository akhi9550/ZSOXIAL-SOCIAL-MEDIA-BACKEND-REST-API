package interfaces

import (
	"mime/multipart"

	"github.com/akhi9550/api-gateway/pkg/utils/models"
)

type PostClient interface {
	CreatePost(userID int, data models.PostRequest, file []*multipart.FileHeader, users models.Tags) (models.PostResponse, error)
	GetPost(userID int, postID int) (models.PostResponse, error)
	UpdatePost(userID int, data models.UpdatePostReq) (models.PostResponse, error)
	DeletePost(userID int, postID int) error
	// GetAllPost(userID int) ([]models.PostResponse, error)
	// LikePost(userID int, postID int) (models.LikePostResponse, error)
	// PostComment(userID int, data models.PostCommentReq) (models.PostCommentResponse, error)
	// GetComment(postID int) (models.GetCommentResponse, error)
	// ReplyComment(userID int, data models.CommentReply) error
	// DeleteComment(userID int, commentID int) error
}
