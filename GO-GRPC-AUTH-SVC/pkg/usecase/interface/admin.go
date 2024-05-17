package interfaces

import "github.com/akhi9550/auth-svc/pkg/utils/models"

type AdminUseCase interface {
	AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error) 
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	AdminBlockUser(id uint) error
	AdminUnBlockUser(id uint) error
	ShowUserReports(page, count int) ([]models.UserReports, error)
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.PostResponse, error)
	RemovePost(postID int) error
}
