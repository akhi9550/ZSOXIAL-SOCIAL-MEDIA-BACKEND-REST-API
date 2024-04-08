package usecase

import (
	"errors"

	authclientinterfaces "github.com/akhi9550/post-svc/pkg/client/interface"
	"github.com/akhi9550/post-svc/pkg/helper"
	interfaces "github.com/akhi9550/post-svc/pkg/repository/interface"
	services "github.com/akhi9550/post-svc/pkg/usecase/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"github.com/google/uuid"
)

type storyUseCase struct {
	storyRepository interfaces.StoryRepository
	authClient      authclientinterfaces.NewauthClient
}

func NewStoryUseCase(repository interfaces.StoryRepository, authclient authclientinterfaces.NewauthClient) services.StoryUseCase {
	return &storyUseCase{
		storyRepository: repository,
		authClient:      authclient,
	}

}

func (s *storyUseCase) CreateStory(userID int, file []byte) (models.CreateStoryResponse, error) {
	userExist := s.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return models.CreateStoryResponse{}, errors.New("user doesn't exist")
	}
	fileUID := uuid.New()
	fileName := fileUID.String()
	s3Path := "story/" + fileName + ".jpg"
	url, err := helper.AddImageToAwsS3(file, s3Path)
	if err != nil {
		return models.CreateStoryResponse{}, err
	}
	data, err := s.storyRepository.CreateStory(userID, url)
	if !userExist {
		return models.CreateStoryResponse{}, err
	}
	userData, err := s.authClient.UserData(int(userID))
	if err != nil {
		return models.CreateStoryResponse{}, err
	}
	return models.CreateStoryResponse{
		Author:    userData,
		Story:     data.Story,
		CreatedAt: data.CreatedAt,
	}, nil
}

func (s *storyUseCase) GetStory(userID, viewer int) ([]models.CreateStoryResponses, error) {
	userExist := s.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return []models.CreateStoryResponses{}, errors.New("user doesn't exist")
	}
	data, err := s.storyRepository.GetStory(userID,viewer)
	if err != nil {
		return []models.CreateStoryResponses{}, err
	}
	var storyResponses []models.CreateStoryResponses
	for _, singleStory := range data {
		userData, err := s.authClient.UserData(int(userID))
		if err != nil {
			return nil, err
		}
		storyResponses = append(storyResponses, models.CreateStoryResponses{
			Author:    userData,
			StoryID:   singleStory.StoryID,
			Story:     singleStory.Story,
			CreatedAt: singleStory.CreatedAt,
		})
	}
	return storyResponses, nil
}

func (s *storyUseCase) DeleteStory(userID, storyID int) error {
	userExist := s.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := s.storyRepository.CheckStoryAvalilabilityWithID(userID, storyID)
	if !ok {
		return errors.New("story doesn't exist")
	}
	err := s.storyRepository.DeleteStory(userID, storyID)
	if err != nil {
		return err
	}
	return nil
}

func (s *storyUseCase) LikeStory(userID, storyID int) error {
	userExist := s.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := s.storyRepository.CheckStoryAvalilabilityWithID(userID, storyID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	ok = s.storyRepository.CheckAlreadyLiked(userID, storyID)
	if ok {
		return errors.New("already liked")
	}
	err := s.storyRepository.LikeStory(userID, storyID)
	if err != nil {
		return err
	}
	return nil
}

func (s *storyUseCase) UnLikeStory(userID, storyID int) error {
	userExist := s.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := s.storyRepository.CheckStoryAvalilabilityWithID(userID, storyID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	ok = s.storyRepository.CheckAlreadyLiked(userID, storyID)
	if !ok {
		return errors.New("already unliked")
	}
	// ok = s.storyRepository.CheckAlreadyLiked(userID, storyID)
	// if !ok {
	// 	s.storyRepository.LikeStory(userID, storyID)
	// 	return errors.New("")
	// }
	err := s.storyRepository.UnLikeStory(userID, storyID)
	if err != nil {
		return err
	}
	return nil
}
