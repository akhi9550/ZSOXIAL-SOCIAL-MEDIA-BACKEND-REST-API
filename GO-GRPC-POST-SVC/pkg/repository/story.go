package repository

import (
	"time"

	interfaces "github.com/akhi9550/post-svc/pkg/repository/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"gorm.io/gorm"
)

type storyRepository struct {
	DB *gorm.DB
}

func NewStoryRepository(DB *gorm.DB) interfaces.StoryRepository {
	return &storyRepository{
		DB: DB,
	}
}

func (s *storyRepository) CreateStory(userID int, file string) (models.CreateStory, error) {
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour * 24)
	var a models.CreateStories
	err := s.DB.Raw(`INSERT INTO stories (user_id,url,start_time,end_time) VALUES (?,?,?,?) RETURNING id,url,start_time`, userID, file, startTime, endTime).Scan(&a).Error
	if err != nil {
		return models.CreateStory{}, err
	}
	err = s.DB.Exec(`INSERT INTO archive_stories (user_id,story_id) VALUES($1,$2)`, userID, a.ID).Error
	if err != nil {
		return models.CreateStory{}, err
	}
	return models.CreateStory{
		Story:     a.Url,
		CreatedAt: a.StartTime,
	}, nil

}

func (s *storyRepository) GetStory(userID, viewer int) ([]models.CreateStoriesResponse, error) {
	var response []models.CreateStoriesResponse
	err := s.DB.Raw(`SELECT id,url,start_time FROM stories WHERE user_id = ? AND is_valid = 'false'`, userID).Scan(&response).Error
	if err != nil {
		return []models.CreateStoriesResponse{}, err
	}
	for _, i := range response {
		err = s.DB.Exec(`INSERT INTO view_stories (story_id,viewer_id) VALUES (?,?)`, i.StoryID, viewer).Error
		if err != nil {
			return []models.CreateStoriesResponse{}, err
		}
	}
	return response, nil
}

func (s *storyRepository) CheckStoryAvalilabilityWithID(userID, storyID int) bool {
	var count int
	err := s.DB.Raw(`SELECT COUNT(*) FROM stories WHERE id = ? AND user_id = ?`, storyID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (s *storyRepository) DeleteStory(userID, storyID int) error {
	err := s.DB.Exec(`DELETE FROM stories WHERE id = ? AND user_id = ?`, storyID, userID).Error
	if err != nil {
		return err
	}
	err = s.DB.Exec(`DELETE FROM archive_stories WHERE story_id = ? AND user_id = ?`, storyID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storyRepository) CheckAlreadyLiked(userID, storyID int) bool {
	var count int
	err := s.DB.Raw(`SELECT COUNT(*) FROM story_likes WHERE story_id = ?  AND liked_user = ?`, storyID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (s *storyRepository) LikeStory(userID, storyID int) error {
	err := s.DB.Exec(`INSERT INTO story_likes (story_id, liked_user, created_at) VALUES ($1, $2, NOW())`, storyID, userID).Error
	if err != nil {
		return err
	}
	err = s.DB.Exec(`UPDATE stories SET likes_count = likes_count + 1 WHERE id = ?`, storyID).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storyRepository) PostedStoryUser(storyID int) (int, error) {
	var id int
	err := s.DB.Raw(`SELECT user_id FROM stories WHERE id = ?`, storyID).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *storyRepository) UnLikeStory(userID, storyID int) error {
	err := s.DB.Exec(`UPDATE stories SET likes_count = likes_count - 1 WHERE id = ?`, storyID).Error
	if err != nil {
		return err
	}
	err = s.DB.Exec(`DELETE FROM story_likes WHERE liked_user = ? AND story_id = ?`, userID, storyID).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storyRepository) ViewersDetails(storyID int) ([]models.Viewer, error) {
	var response []models.Viewer
	err := s.DB.Raw(`SELECT viewer_id FROM view_stories WHERE story_id = ?`, storyID).Scan(&response).Error
	if err != nil {
		return []models.Viewer{}, err
	}
	return response, nil
}
func (s *storyRepository) LikedUser(storyID int) ([]models.Likeduser, error) {
	var response []models.Likeduser
	err := s.DB.Raw(`SELECT liked_user FROM story_likes WHERE story_id = ?`, storyID).Scan(&response).Error
	if err != nil {
		return []models.Likeduser{}, err
	}
	return response, nil
}
