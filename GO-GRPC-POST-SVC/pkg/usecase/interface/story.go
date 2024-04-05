package interfaces

import "github.com/akhi9550/post-svc/pkg/utils/models"

type StoryUseCase interface {
	CreateStory(userID int, file []byte) (models.CreateStoryResponse, error)
	GetStory(userID int) ([]models.CreateStoryResponses, error)
	DeleteStory(userID, storyID int) error
	LikeStory(userID, storyID int) error
	UnLikeStory(userID, storyID int) error
}
