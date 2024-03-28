package interfaces

import (
	"github.com/akhi9550/auth-svc/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(user models.UserSignUpRequest, file []byte) (*models.ReponseWithToken, error)
	UserLogin(user models.UserLoginRequest) (*models.ReponseWithToken, error)
	ForgotPassword(phone string) error
	ForgotPasswordVerifyAndChange(model models.ForgotVerify) error
	UserDetails(userID int) (models.UsersProfileDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetail,file []byte, userID int) (models.UsersProfileDetails, error)
	ChangePassword(id int, change models.ChangePassword) error
}
