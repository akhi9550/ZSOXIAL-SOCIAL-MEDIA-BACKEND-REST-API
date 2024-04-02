package interfaces

import (
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

type PostRepository interface {
	CheckUserAvalilabilityWithUserID(userID int) bool
	CheckMediaAvalilabilityWithID(typeID int) bool
	CheckPostAvalilabilityWithID(postID int) bool
	UserData(userID int) (models.UserData, error)
	CreatePost(userID int, Caption string, TypeId int, file string, users []models.Tag) (models.Response, []models.Tag, error)
	GetPost(userID, postID int) (models.Response, []models.Tag, error)
	UpdateCaption(postID, userID int, caption string) error
	UpdateTypeID(userID, PostID, TypeID int) error
	UpdateTags(userID, PostID int, tag []models.Tag) error
	PostDetails(PostID, userID int) (models.Response, []models.Tag, error)
	DeletePost(userID, postID int) error
}
