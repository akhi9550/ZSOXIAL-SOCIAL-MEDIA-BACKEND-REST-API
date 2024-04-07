package interfaces

import (
	"github.com/akhi9550/auth-svc/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(user models.UserSignUpRequest) (*models.ReponseWithToken, error)
	UserLogin(user models.UserLoginRequest) (*models.ReponseWithToken, error)
	ForgotPassword(phone string) error
	ForgotPasswordVerifyAndChange(model models.ForgotVerify) error
	UserDetails(userID int) (models.UsersProfileDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetail, file []byte, userID int) (models.UsersProfileDetails, error)
	ChangePassword(id int, change models.ChangePassword) error
	CheckUserAvalilabilityWithUserID(userID int) (bool, error)
	UserData(userID int) (models.UserData, error)
	CheckUserAvalilabilityWithTagUserID(users []models.Tag) (bool, error)
	GetUserNameWithTagUserID(users []models.Tag) ([]models.UserTag, error)
	ReportUser(userID int, req models.ReportRequest) error
}
