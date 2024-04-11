package interfaces

import (
	"mime/multipart"

	"github.com/akhi9550/api-gateway/pkg/utils/models"
)

type PostClient interface {
	CreatePost(userID int, data models.PostRequest, file *multipart.FileHeader, users []string) (models.PostResponse, error)
	GetPost(userID int, postID int) (models.PostResponse, error)
	UpdatePost(userID int, data models.UpdatePostReq, user []string) (models.UpdateResponse, error)
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
	ShowAllPostComments(PostID int) ([]models.AllCommentsAndReplies, error)
	ReportPost(userID int, req models.ReportPostRequest) error
	ReplyComment(userID int, req models.ReplyCommentReq) (models.ReplyReposne, error)
	SavedPost(userID, postID int) error
	UnSavedPost(userID, postID int) error
	GetSavedPost(userID int) ([]models.PostResponse, error)
	CreateStory(userID int, file *multipart.FileHeader) (models.CreateStoryResponse, error)
	GetStory(userID, viewer int) ([]models.CreateStoryResponses, error)
	DeleteStory(userID, storyID int) error
	LikeStory(userID, storyID int) error
	UnLikeStory(userID, storyID int) error
	StoryDetails(userID, storyID int) (models.StoryDetails,error)
}
