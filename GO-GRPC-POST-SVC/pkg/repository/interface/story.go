package interfaces

import (
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

type StoryRepository interface {
	CreateStory(userID int, file string) (models.CreateStory, error)
	GetStory(userID, viewer int) ([]models.CreateStoriesResponse, error)
	CheckStoryAvalilabilityWithID(userID, storyID int) bool
	DeleteStory(userID, storyID int) error
	CheckAlreadyLiked(userID, storyID int) bool
	LikeStory(userID, storyID int) error
	PostedStoryUser(storyID int) (int, error)
	UnLikeStory(userID, storyID int) error
	ViewersDetails(storyID int) ([]models.Viewer, error)
	LikedUser(storyID int) ([]models.Likeduser, error)
}
