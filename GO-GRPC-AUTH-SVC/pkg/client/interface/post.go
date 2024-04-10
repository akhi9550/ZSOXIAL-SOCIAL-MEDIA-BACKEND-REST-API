package interfaces

import "github.com/akhi9550/auth-svc/pkg/utils/models"

type NewPostClient interface {
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.PostResponse, error)
	CheckPostIDByID(postID int) bool
	RemovePost(postID int) error
}
