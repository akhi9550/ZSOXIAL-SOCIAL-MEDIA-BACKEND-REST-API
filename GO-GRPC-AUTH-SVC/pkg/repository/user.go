package repository

import (
	"errors"

	"github.com/akhi9550/auth-svc/pkg/domain"
	interfaces "github.com/akhi9550/auth-svc/pkg/repository/interface"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		DB: DB,
	}
}

func (ur *userRepository) CheckUserExistsByUsername(username string) (*domain.User, error) {
	var user domain.User
	res := ur.DB.Where(&domain.User{Username: username}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (ur *userRepository) CheckUserExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	res := ur.DB.Where(&domain.User{Email: email}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (ur *userRepository) CheckUserExistsByPhone(phone string) (*domain.User, error) {
	var user domain.User
	res := ur.DB.Where(&domain.User{Phone: phone}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (ur *userRepository) UserSignUp(user models.UserSignUpRequest, url string) (models.UserResponse, error) {
	var signupDetails models.UserResponse
	err := ur.DB.Raw(`INSERT INTO users (firstname,lastname,username,dob,gender,phone,email,password,bio,imageurl,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,NOW()) RETURNING id,username,imageurl`, user.Firstname, user.Lastname, user.Username, user.Dob, user.Gender, user.Phone, user.Email, user.Password, user.Bio, url).Scan(&signupDetails).Error
	if err != nil {
		return models.UserResponse{}, err
	}
	return signupDetails, nil
}

func (ur *userRepository) FindUserByEmail(user models.UserLoginRequest) (models.UserResponsewithPassword, error) {
	var userDetails models.UserResponsewithPassword
	err := ur.DB.Raw("SELECT * FROM users WHERE email=? and blocked=false and isadmin=false", user.Email).Scan(&userDetails).Error
	if err != nil {
		return models.UserResponsewithPassword{}, errors.New("error checking user details")
	}
	return userDetails, nil
}
func (ur *userRepository) FindUserBlockorNot(email string) (bool, error) {
	var a bool
	if err := ur.DB.Raw("SELECT blocked FROM users WHERE email = ?", email).Scan(&a).Error; err != nil {
		return false, err
	}
	return a, nil
}

func (ur *userRepository) FindUserByMobileNumber(phone string) bool {
	var count int
	if err := ur.DB.Raw("SELECT count(*) FROM users WHERE phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0

}

func (ur *userRepository) FindIdFromPhone(phone string) (int, error) {
	var id int
	if err := ur.DB.Raw("SELECT id FROM users WHERE phone=?", phone).Scan(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

func (ur *userRepository) ChangePassword(id int, password string) error {
	err := ur.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", password, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	err := ur.DB.Raw("SELECT firstname,lastname,username,dob,gender,phone,email,bio,imageurl FROM users WHERE id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Username, &userDetails.Dob, &userDetails.Gender, &userDetails.Phone, &userDetails.Email, &userDetails.Bio, &userDetails.Imageurl)
	if err != nil {
		return models.UsersProfileDetails{}, errors.New("could not get user details")
	}
	return userDetails, nil
}

func (ur *userRepository) CheckUserAvailabilityWithUserID(userID int) bool {
	var count int
	if err := ur.DB.Raw("SELECT COUNT(*) FROM users WHERE id= ?", userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (ur *userRepository) UpdateFirstName(firstname string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET firstname= ? WHERE id = ?", firstname, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateLastName(lastname string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET lastname= ? WHERE id = ?", lastname, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserName(username string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET username= ? WHERE id = ?", username, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateDOB(dob string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET dob= ? WHERE id = ?", dob, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateGender(gender string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET gender= ? WHERE id = ?", gender, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserPhone(phone string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET phone= ? WHERE id = ?", phone, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserEmail(email string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET email= ? WHERE id = ?", email, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateBIO(bio string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET bio= ? WHERE id = ?", bio, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdatePhoto(imageurl string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET imageurl= ? WHERE id = ?", imageurl, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetPassword(id int) (string, error) {
	var userPassword string
	err := ur.DB.Raw("SELECT password FROM users WHERE id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil
}

func (ur *userRepository) ExistUsername(username string) bool {
	var count int
	if err := ur.DB.Raw("SELECT count(*) FROM users WHERE username = ?", username).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (ur *userRepository) ExistPhone(phone string) bool {
	var count int
	if err := ur.DB.Raw("SELECT count(*) FROM users WHERE phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (ur *userRepository) ExistEmail(email string) bool {
	var count int
	if err := ur.DB.Raw("SELECT count(*) FROM users WHERE email = ?", email).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
