package interfaces

import (
	"github.com/akhi9550/auth-svc/pkg/utils/models"
	"github.com/akhi9550/auth-svc/pkg/domain"
)

type AdminRepository interface {
	CheckAdminExistsByEmail(email string) (*domain.User, error)
	FindAdminByEmail(user models.AdminLoginRequest) (models.AdminResponsewithPassword, error)
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	GetUserByID(id uint) (domain.User, error)
	AdminBlockUserByID(user domain.User) error
}
