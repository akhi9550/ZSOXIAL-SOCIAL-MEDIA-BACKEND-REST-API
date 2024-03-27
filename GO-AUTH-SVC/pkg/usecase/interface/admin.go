package interfaces

import "github.com/akhi9550/pkg/utils/models"

type AdminUseCase interface {
	AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error) 
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	AdminBlockUser(id uint) error
	AdminUnBlockUser(id uint) error
}
