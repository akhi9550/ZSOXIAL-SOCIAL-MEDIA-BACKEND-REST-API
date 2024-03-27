package interfaces

import (
	"github.com/akhi9550/auth-svc/pkg/domain"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
)

type OtpRepository interface {
	FindUserByPhoneNumber(phone string) (*domain.User, error)
	UserDetailsUsingPhone(phone string) (models.UserResponse, error)
	FindUsersByEmail(email string) (bool, error)
	GetUserPhoneByEmail(email string) (string, error)
}
