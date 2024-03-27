package repository

import (
	"errors"
	"fmt"

	"github.com/akhi9550/pkg/domain"
	interfaces "github.com/akhi9550/pkg/repository/interface"
	"github.com/akhi9550/pkg/utils/models"
	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}

func (ad *adminRepository) CheckAdminExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	res := ad.DB.Where(&domain.User{Email: email}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (ad *adminRepository) FindAdminByEmail(user models.AdminLoginRequest) (models.AdminResponsewithPassword, error) {
	var userDetails models.AdminResponsewithPassword
	err := ad.DB.Raw("SELECT * FROM users WHERE email=? and isadmin=true", user.Email).Scan(&userDetails).Error
	if err != nil {
		return models.AdminResponsewithPassword{}, errors.New("error checking user details")
	}
	return userDetails, nil
}

func (ad *adminRepository) ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error) {
	var user []models.UserDetailsAtAdmin
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	err := ad.DB.Raw("SELECT id,firstname,lastname,username,dob,gender,phone,email,imageurl,created_at,blocked FROM users WHERE isadmin='false' limit ? offset ?", count, offset).Scan(&user).Error
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return user, nil
}
func (ad *adminRepository) GetUserByID(userID uint) (domain.User, error) {
	var count int
	if err := ad.DB.Raw("SELECT COUNT(*) FROM users WHERE id=?", userID).Scan(&count).Error; err != nil {

		return domain.User{}, err
	}
	if count < 1 {
		return domain.User{}, errors.New("user for the given id does not exist")

	}
	var userDetails domain.User
	if err := ad.DB.Raw("SELECT * FROM users WHERE id=?", userID).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}
	return userDetails, nil
}

func (ad *adminRepository) AdminBlockUserByID(user domain.User) error {
	err := ad.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}
