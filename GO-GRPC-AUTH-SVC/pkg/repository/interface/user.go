package interfaces

import (
	"github.com/akhi9550/auth-svc/pkg/utils/models"
	"github.com/akhi9550/auth-svc/pkg/domain"
)

type UserRepository interface {
	CheckUserExistsByUsername(username string) (*domain.User, error)
	CheckUserExistsByEmail(email string) (*domain.User, error)
	CheckUserExistsByPhone(phone string) (*domain.User, error)
	UserSignUp(user models.UserSignUpRequest,url string) (models.UserResponse, error)
	FindUserByEmail(user models.UserLoginRequest) (models.UserResponsewithPassword, error)
	FindUserBlockorNot(email string)(bool,error)
	FindUserByMobileNumber(phone string) bool
	FindIdFromPhone(phone string) (string, error)
	ChangePassword(phone string, password string) error
	UserDetails(userID int) (models.UsersProfileDetails, error)
	CheckUserAvailabilityWithUserID(userID int) bool
	UpdateFirstName(firstname string, userID int) error
	UpdateLastName(lastname string, userID int) error
	UpdateUserName(username string, userID int) error
	UpdateDOB(dob string, userID int) error
	UpdateGender(gender string, userID int) error
	UpdateUserPhone(phone string, userID int) error
	UpdateUserEmail(email string, userID int) error
	UpdateBIO(bio string, userID int) error
	UpdatePhoto(image string, userID int) error
	GetPassword(id int) (string, error)
	ExistUsername(username string) bool
	ExistPhone(phone string) bool
	ExistEmail(email string) bool
	Changepassword(phone int, password string) error
}
