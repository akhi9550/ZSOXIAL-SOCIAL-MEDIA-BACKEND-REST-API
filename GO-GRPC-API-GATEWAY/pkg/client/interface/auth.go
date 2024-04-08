package interfaces

import (
	"mime/multipart"

	"github.com/akhi9550/api-gateway/pkg/utils/models"
)

type AuthClient interface {
	UserSignUp(user models.UserSignUpRequest) (*models.ReponseWithToken, error)
	UserLogin(user models.UserLoginRequest) (*models.ReponseWithToken, error)
	ForgotPassword(phone string) error
	ForgotPasswordVerifyAndChange(model models.ForgotVerify) error
	UserDetails(userID int) (models.UsersProfileDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetail, file *multipart.FileHeader, userID int) (models.UsersProfileDetails, error)
	ChangePassword(userID int, change models.ChangePassword) error
	SendOtp(phone string) error
	VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error)
	AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error)
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	AdminBlockUser(userID int) error
	AdminUnblockUser(userID int) error
	ReportUser(userID int, req models.ReportRequest) error
	FollowREQ(userID, FollowingID int) error
	ShowFollowREQ(userID int) ([]models.FollowingRequests, error)
	AcceptFollowREQ(userID, FollowingID int) error
	Following(userID int) ([]models.FollowingResponse, error)
	Follower(userID int) ([]models.FollowingResponse, error)
}
