package interfaces

import (
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

type StoryRepository interface {
	CreateStory(userID int, file string) (models.CreateStory, error)
	GetStory(userID int) ([]models.CreateStoriesResponse, error)
	CheckStoryAvalilabilityWithID(userID, storyID int) bool
	DeleteStory(userID, storyID int) error
	CheckAlreadyLiked(userID, storyID int) bool
	LikeStory(userID, storyID int) error
	UnLikeStory(userID, storyID int) error
}
