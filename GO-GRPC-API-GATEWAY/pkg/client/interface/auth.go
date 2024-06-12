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
	SpecificUserDetails(userID int) (models.UsersDetails, error)
	UserDetails(userID int) (models.UsersDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetail, file *multipart.FileHeader, userID int) (models.UsersProfileDetails, error)
	ChangePassword(userID int, change models.ChangePassword) error
	SendOtp(phone string) error
	VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error)
	ReportUser(userID int, req models.ReportRequest) error
	FollowREQ(userID, FollowingID int) error
	ShowFollowREQ(userID int) ([]models.FollowingRequests, error)
	AcceptFollowREQ(userID, FollowUserID int) error
	UnFollow(userID, UnFollowUserID int) error
	Following(userID int) ([]models.FollowingResponse, error)
	Follower(userID int) ([]models.FollowingResponse, error)
	AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error)
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	AdminBlockUser(userID int) error
	AdminUnblockUser(userID int) error
	ShowUserReports(page, count int) ([]models.UserReports, error)
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.PostResponse, error)
	RemovePost(postID int) error
	SearchUser(req models.SearchUser) ([]models.SearchResult, error)
	VideoCallKey(userID, oppositeUser int) (string, error)
	CreateGroup(userID int,req models.GroupReq, users []string,file *multipart.FileHeader)error
	ExitFormGroup(userID,groupID int)error
	ShowGroups(userID int)([]models.Groups,error)
	ShowGroupMembers(userID,groupID int)([]models.Mebmers,error)
}
