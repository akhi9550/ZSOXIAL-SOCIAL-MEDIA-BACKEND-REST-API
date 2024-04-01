package interfaces

import "github.com/akhi9550/post-svc/pkg/utils/models"

type StoryUseCase interface {
	CreateStory(userID uint, data models.StoryUrl) (models.StoryResponse, error)
	DeleteStory(userID uint, storyID uint) error
	Likestory(userID uint, storyID uint) error
}
