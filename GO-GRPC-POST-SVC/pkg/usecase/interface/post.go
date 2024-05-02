package interfaces

import (
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

type PostUseCase interface {
	CreatePost(userID int, data models.PostRequest, file []byte, users []models.Tag) (models.PostResponse, error)
	GetPost(userID int, postID int) (models.PostResponse, error)
	UpdatePost(userID int, data models.UpdatePostReq, tag []models.Tag) (models.UpdateResponse, error)
	DeletePost(userID int, postID int) error
	GetAllPost(userID int) ([]models.PostResponse, error)
	ArchivePost(userID, PostID int) error
	UnArchivePost(userID, PostID int) error
	GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error)
	LikePost(userID int, postID int) (models.LikePostResponse, error)
	UnLinkPost(userID int, postID int) error
	PostComment(userID int, data models.PostCommentReq) (models.PostComment, error)
	DeleteComment(userID, CommentID int) error
	GetAllPostComments(PostID int) ([]models.PostCommentResponse, error)
	ReplyComment(userID int, req models.ReplyCommentReq) (models.ReplyReposne, error)
	ShowAllPostComments(PostID int) ([]models.AllCommentsAndReplies, error)
	ReportPost(userID int, req models.ReportRequest) error
	SavedPost(userID, postID int) error
	UnSavedPost(userID, postID int) error
	GetSavedPost(userID int) ([]models.PostResponse, error)
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.PostResponse, error)
	CheckPostIDByID(postID int) bool
	RemovePost(postID int) error
	Home(userID int) ([]models.PostResponse, error)
}
