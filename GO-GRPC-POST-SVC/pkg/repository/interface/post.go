package interfaces

import (
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

type PostRepository interface {
	CheckUserAvalilabilityWithUserID(userID int) bool
	CheckMediaAvalilabilityWithID(typeID int) bool
	CheckPostAvalilabilityWithID(postID int) bool
	CheckPostedUserID(userID, PostID int)bool
	UserData(userID int) (models.UserData, error)
	CreatePost(userID int, Caption string, TypeId int, file string, users []models.Tag) (models.Response, []models.Tag, error)
	GetPost(postID int) (models.Responses, error)
	GetTagUser(postID int) ([]models.Tag, error)
	UpdateCaption(postID, userID int, caption string) error
	UpdateTypeID(userID, PostID, TypeID int) error
	UpdateTags(userID, PostID int, tag []models.Tag) error
	PostDetails(PostID, userID int) (models.Response, []models.Tag, error)
	DeletePost(userID, postID int) error
	GetPostAll(userID int) ([]models.Response, error)
	CheckPostAvalilabilityWithUserID(postID, userID int) bool
	ArchivePost(userID, postID int) error
	UnArchivePost(userID, postID int) error
	GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error)
	CheckAlreadyLiked(userID, PostID int) bool
	LikePost(userID, postID int) (models.LikesReponse, error)
	GetPostedUserID(postID int) (int, error)
	UnLikePost(userID, postID int) error
	PostComment(userID int, data models.PostCommentReq) (models.PostComments, error)
	CheckUserWithUserID(userID int) bool
	CheckCommentWithID(CommentID int) bool
	DeleteComment(userID, CommentID int) error
	GetAllPostComments(PostID int) ([]models.PostCommentResponses, error)
	GetCommentsByPostID(postID int) ([]models.AllComments, error)
	GetRepliesByID(PostID, CommentID int) ([]models.Replies, error)
	AlreadyReported(userID, postID int) bool
	ReportPost(userID int, req models.ReportRequest) error
	AllReadyExistReply(userID, CommentID int) bool
	ReplyComment(userID int, req models.ReplyCommentReq) (models.PostComments, models.ReplyResponse, error)
	SavedPost(userID, postID int) error
	AllReadyExistPost(userID, postID int) bool
	UnSavedPost(userID, postID int) error
	GetSavedPost(userID int) ([]models.SavedResponse, error)
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.Responses, error)
	CheckPostIDByID(postID int) bool
	RemovePost(postID int) error
	Home(users []models.Users) ([]models.Responses, error)
}
