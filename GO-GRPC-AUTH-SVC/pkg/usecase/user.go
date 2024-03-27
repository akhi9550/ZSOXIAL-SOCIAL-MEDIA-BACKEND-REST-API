package usecase

import (
	"errors"

	"github.com/akhi9550/auth-svc/pkg/config"
	"github.com/akhi9550/auth-svc/pkg/helper"
	interfaces "github.com/akhi9550/auth-svc/pkg/repository/interface"
	services "github.com/akhi9550/auth-svc/pkg/usecase/interface"
	"github.com/akhi9550/auth-svc/pkg/utils/models"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepository interfaces.UserRepository
}

func NewUserUseCase(repository interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepository: repository,
	}
}
func (ur *userUseCase) UserSignUp(user models.UserSignUpRequest) (*models.ReponseWithToken, error) {
	username, err := ur.userRepository.CheckUserExistsByUsername(user.Username)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("error with server")
	}
	if username != nil {
		return &models.ReponseWithToken{}, errors.New("user with this username is already exists")
	}
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("error with server")
	}
	if email != nil {
		return &models.ReponseWithToken{}, errors.New("user with this email is already exists")
	}

	phone, err := ur.userRepository.CheckUserExistsByPhone(user.Phone)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("error with server")
	}
	if phone != nil {
		return &models.ReponseWithToken{}, errors.New("user with this phone is already exists")
	}

	hashPassword, err := helper.PasswordHash(user.Password)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("error in hashing password")
	}
	user.Password = hashPassword
	userData, err := ur.userRepository.UserSignUp(user)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("could not add the user")
	}
	accessToken, err := helper.GenerateAccessTokenUser(userData)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("couldn't create access token due to error")
	}
	RefreshToken, err := helper.GenerateRefreshTokenUser(userData)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("couldn't create access token due to error")
	}
	return &models.ReponseWithToken{
		Users:        userData,
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
	}, nil
}

func (ur *userUseCase) UserLogin(user models.UserLoginRequest) (*models.ReponseWithToken, error) {
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("error with server")
	}
	if email == nil {
		return &models.ReponseWithToken{}, errors.New("email doesn't exist")
	}
	userdeatils, err := ur.userRepository.FindUserByEmail(user)
	if err != nil {
		return &models.ReponseWithToken{}, err
	}
	ok, err := ur.userRepository.FindUserBlockorNot(user.Email)
	if ok {
		return &models.ReponseWithToken{}, errors.New("user is blocked")
	}
	if err != nil {
		return &models.ReponseWithToken{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userdeatils.Password), []byte(user.Password))
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("password not matching")
	}
	var userData models.UserResponse
	err = copier.Copy(&userData, &userdeatils)
	if err != nil {
		return &models.ReponseWithToken{}, err
	}
	accessToken, err := helper.GenerateAccessTokenUser(userData)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("couldn't create access token due to error")
	}
	RefreshToken, err := helper.GenerateRefreshTokenUser(userData)
	if err != nil {
		return &models.ReponseWithToken{}, errors.New("couldn't create access token due to error")
	}
	return &models.ReponseWithToken{
		Users:        userData,
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
	}, nil
}

func (ur *userUseCase) ForgotPasswordSend(phone string) error {
	cfg, _ := config.LoadConfig()
	ok := ur.userRepository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	_, err := helper.TwilioSendOTP(phone, cfg.SERVICESSID)
	if err != nil {
		return errors.New("error ocurred while generating OTP")
	}
	return nil
}

func (ur *userUseCase) ForgotPasswordVerifyAndChange(model models.ForgotVerify) error {
	cfg, _ := config.LoadConfig()
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(cfg.SERVICESSID, model.Otp, model.Phone)
	if err != nil {
		return errors.New("error while verifying")
	}

	id, err := ur.userRepository.FindIdFromPhone(model.Phone)
	if err != nil {
		return errors.New("cannot find user from mobile number")
	}

	newpassword, err := helper.PasswordHashing(model.NewPassword)
	if err != nil {
		return errors.New("error in hashing password")
	}

	if err := ur.userRepository.ChangePassword(id, string(newpassword)); err != nil {
		return errors.New("could not change password")
	}

	return nil
}

func (ur *userUseCase) UserDetails(userID int) (models.UsersProfileDetails, error) {
	return ur.userRepository.UserDetails(userID)
}

func (ur *userUseCase) UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error) {
	userExist := ur.userRepository.CheckUserAvailabilityWithUserID(userID)
	if !userExist {
		return models.UsersProfileDetails{}, errors.New("user doesn't exist")
	}
	if userDetails.Firstname != "" {
		ur.userRepository.UpdateFirstName(userDetails.Firstname, userID)
	}
	if userDetails.Lastname != "" {
		ur.userRepository.UpdateLastName(userDetails.Lastname, userID)
	}
	if userDetails.Lastname != "" {
		ok := ur.userRepository.ExistUsername(userDetails.Username)
		if ok {
			return models.UsersProfileDetails{}, errors.New("username already exist")
		}
		ur.userRepository.UpdateUserName(userDetails.Username, userID)
	}
	if userDetails.Lastname != "" {
		ur.userRepository.UpdateDOB(userDetails.Dob, userID)
	}
	if userDetails.Lastname != "" {
		ur.userRepository.UpdateGender(userDetails.Gender, userID)
	}
	if userDetails.Phone != "" {
		ok := ur.userRepository.ExistPhone(userDetails.Phone)
		if ok {
			return models.UsersProfileDetails{}, errors.New("phone already exist")
		}
		ur.userRepository.UpdateUserPhone(userDetails.Phone, userID)
	}
	if userDetails.Email != "" {
		ok := ur.userRepository.ExistEmail(userDetails.Email)
		if ok {
			return models.UsersProfileDetails{}, errors.New("email already exist")
		}
		ur.userRepository.UpdateUserEmail(userDetails.Email, userID)
	}
	if userDetails.Phone != "" {
		ur.userRepository.UpdateBIO(userDetails.Bio, userID)
	}
	if userDetails.Phone != "" {
		ur.userRepository.UpdatePhoto(userDetails.Imageurl, userID)
	}
	return ur.userRepository.UserDetails(userID)
}

func (ur *userUseCase) ChangePassword(id int, change models.ChangePassword) error {
	userPassword, err := ur.userRepository.GetPassword(id)
	if err != nil {
		return errors.New("internal error")
	}
	err = helper.CompareHashAndPassword(userPassword, change.Oldpassword)
	if err != nil {
		return errors.New("password incorrect")
	}
	if change.Password != change.Repassword {
		return errors.New("password doesn't match")
	}
	newpassword, err := helper.PasswordHash(change.Password)
	if err != nil {
		return errors.New("error in hashing password")
	}
	return ur.userRepository.ChangePassword(id, string(newpassword))
}
