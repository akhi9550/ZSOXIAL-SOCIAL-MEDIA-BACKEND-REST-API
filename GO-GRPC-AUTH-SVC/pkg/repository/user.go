package repository

import (
	"errors"
	"fmt"

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

func (ur *userRepository) UserSignUp(user models.UserSignUpRequest) (models.UserResponse, error) {
	var signupDetails models.UserResponse
	err := ur.DB.Raw(`INSERT INTO users (firstname,lastname,username,phone,email,password,created_at) VALUES(?,?,?,?,?,?,NOW()) RETURNING id,username,email`, user.Firstname, user.Lastname, user.Username, user.Phone, user.Email, user.Password).Scan(&signupDetails).Error
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

func (ur *userRepository) FindIdFromPhone(phone string) (string, error) {
	var id string
	if err := ur.DB.Raw("SELECT id FROM users WHERE phone=?", phone).Scan(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

func (ur *userRepository) ChangePassword(phone string, password string) error {
	err := ur.DB.Exec("UPDATE users SET password = $1 WHERE phone = $2", password, phone).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) SpecificUserDetails(userID int) (models.UsersDetails, error) {
	var userDetails models.UsersDetails
	err := ur.DB.Raw("SELECT firstname, lastname, username, dob, gender, phone, email, bio, imageurl FROM users WHERE id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Username, &userDetails.Dob, &userDetails.Gender, &userDetails.Phone, &userDetails.Email, &userDetails.Bio, &userDetails.Imageurl)
	if err != nil {
		fmt.Println("Error retrieving user details:", err)
		return models.UsersDetails{}, err
	}
	var followerCount int
	result := ur.DB.Raw("SELECT COUNT(*) FROM followers WHERE user_id = ?", userID).Scan(&followerCount)
	if result.Error != nil {
		return models.UsersDetails{}, result.Error
	}
	userDetails.Follower = followerCount

	var followingCount int
	result = ur.DB.Raw("SELECT COUNT(*) FROM followings WHERE user_id = ?", userID).Scan(&followingCount)
	if result.Error != nil {
		return models.UsersDetails{}, result.Error
	}
	userDetails.Following = followingCount
	return userDetails, nil
}

func (ur *userRepository) UserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	err := ur.DB.Raw("SELECT firstname, lastname, username, dob, gender, phone, email, bio, imageurl FROM users WHERE id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Username, &userDetails.Dob, &userDetails.Gender, &userDetails.Phone, &userDetails.Email, &userDetails.Bio, &userDetails.Imageurl)
	if err != nil {
		fmt.Println("Error retrieving user details:", err)
		return models.UsersProfileDetails{}, err
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

func (ur *userRepository) CheckUserAvalilabilityWithUserID(userID int) (bool, error) {
	var count int
	if err := ur.DB.Raw("SELECT COUNT(*) FROM users WHERE id= ?", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *userRepository) UserData(userID int) (models.UserData, error) {
	var user models.UserData
	err := ur.DB.Raw(`SELECT id, username, imageurl FROM users WHERE id = ?`, userID).Scan(&user).Error
	if err != nil {
		return models.UserData{}, err
	}
	return user, nil
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

func (ur *userRepository) Changepassword(id int, password string) error {
	err := ur.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", password, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CheckUserAvalilabilityWithTagUserID(users []models.Tag) (bool, error) {
	var count int
	for _, i := range users {
		if err := ur.DB.Raw("SELECT COUNT(*) FROM users WHERE id = ?", i.User).Scan(&count).Error; err != nil {
			return false, err
		}
	}
	return count > 0, nil
}

func (ur *userRepository) GetUserNameWithTagUserID(users []models.Tag) ([]models.UserTag, error) {
	var data []models.UserTag
	for _, user := range users {
		var userDetails models.UserTag
		if err := ur.DB.Raw("SELECT username FROM users WHERE id = ?", user.User).Scan(&userDetails).Error; err != nil {
			return nil, err
		}
		data = append(data, userDetails)
	}
	return data, nil
}

func (ur *userRepository) GetFollowingUsers(userID int) ([]models.FollowUsers, error) {
	var users []models.FollowUsers
	err := ur.DB.Raw(`SELECT following_user FROM followings WHERE user_id = ?`, userID).Scan(&users).Error
	if err != nil {
		return []models.FollowUsers{}, err
	}
	return users, nil
}

func (ur *userRepository) AlreadyReported(RuserID, userID int) bool {
	var count int
	err := ur.DB.Raw(`SELECT COUNT(*) FROM user_reports WHERE report_user_id = ? AND user_id = ?`, RuserID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (ur *userRepository) ReportUser(userID int, req models.ReportRequest) error {
	err := ur.DB.Exec(`INSERT INTO user_reports (report_user_id,user_id,report) VALUES (?,?,?)`, userID, req.UserID, req.Report).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FollowREQ(userID, FollowingUserID int) error {
	err := ur.DB.Exec(`INSERT INTO following_requests (user_id,following_user,created_at) VALUES(?,?,NOW())`, userID, FollowingUserID).Error
	if err != nil {
		return err
	}
	err = ur.DB.Exec(`INSERT INTO followings (user_id,following_user,created_at) VALUES(?,?,NOW())`, userID, FollowingUserID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) ExistFollowreq(userID, FollowingUserID int) bool {
	var count int
	err := ur.DB.Raw(`SELECT COUNT(*) FROM following_requests WHERE following_user = ? AND user_id = ?`, userID, FollowingUserID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (ur *userRepository) ShowFollowREQ(userID int) ([]models.FollowReqs, error) {
	var response []models.FollowReqs
	err := ur.DB.Raw(`SELECT user_id,created_at FROM following_requests WHERE following_user = ?`, userID).Scan(&response).Error
	if err != nil {
		return []models.FollowReqs{}, err
	}
	return response, nil
}

func (ur *userRepository) CheckRequest(userID, FollowingUserID int) bool {
	var request models.FollowingRequest
	err := ur.DB.Raw(`SELECT following_user, user_id FROM following_requests WHERE following_user = ? AND user_id = ?`, userID, FollowingUserID).Scan(&request).Error
	if err != nil {
		return false
	}
	return request.UserID != 0
}

func (ur *userRepository) AlreadyAccepted(userID, FollowingUserID int) bool {
	var count int
	err := ur.DB.Raw(`SELECT COUNT(*) FROM followers WHERE user_id = ? AND following_user = ?`, userID, FollowingUserID).Scan(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}

	err = ur.DB.Raw(`SELECT COUNT(*) FROM followings WHERE user_id = ? AND following_user = ?`, userID, FollowingUserID).Scan(&count).Error
	if err != nil {
		return false
	}

	return count > 0
}

func (ur *userRepository) AcceptFollowREQ(userID, FollowingUserID int) error {
	err := ur.DB.Exec(`INSERT INTO followers (user_id,following_user,created_at) VALUES(?,?,NOW())`, userID, FollowingUserID).Error
	if err != nil {
		return err
	}

	err = ur.DB.Exec(`DELETE FROM following_requests WHERE user_id = ? AND following_user = ?`, FollowingUserID, userID).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) UnFollow(userID, UnFollowUserID int) error {
	err := ur.DB.Exec(`DELETE FROM followings WHERE user_id = ? AND following_user = ?`, userID, UnFollowUserID).Error
	if err != nil {
		return err
	}
	err = ur.DB.Exec(`DELETE FROM followers WHERE user_id = ? AND following_user = ?`, UnFollowUserID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Following(userID int) ([]models.FollowResp, error) {
	var response []models.FollowResp
	err := ur.DB.Raw(`SELECT following_user FROM followings WHERE user_id = ?`, userID).Scan(&response).Error
	if err != nil {
		return []models.FollowResp{}, err
	}
	return response, nil
}

func (ur *userRepository) Follower(userID int) ([]models.FollowResp, error) {
	var response []models.FollowResp
	err := ur.DB.Raw(`SELECT following_user FROM followers WHERE user_id = ?`, userID).Scan(&response).Error
	if err != nil {
		return []models.FollowResp{}, err
	}
	return response, nil
}

func (ur *userRepository) SearchUser(req models.SearchUser) ([]models.Users, error) {
	var response []models.Users
	// username := "%" + req.Username + "%"
	err := ur.DB.Raw(`SELECT username,imageurl FROM users WHERE username ILIKE '%' || $1 || '%' LIMIT $2 OFFSET $3`, req.Username, req.Limit, req.Offset).Scan(&response).Error
	if err != nil {
		return []models.Users{}, err
	}
	return response, nil
}
